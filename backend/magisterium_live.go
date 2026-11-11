package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"google.golang.org/genai"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow WebSocket connections from allowed origin or any in development
		return true
	},
	ReadBufferSize:  65536,
	WriteBufferSize: 65536,
}

// ClientLiveMessage represents messages sent from the frontend to the backend Live gateway.
type ClientLiveMessage struct {
	Type string `json:"type"` // "realtime_input", "client_content", "audio_stream_end", "ping"
	Data string `json:"data,omitempty"` // Base64 audio PCM (16kHz mono)
	Text string `json:"text,omitempty"` // User prompt text
}

// ServerLiveMessage represents messages sent from the backend Live gateway to the frontend.
type ServerLiveMessage struct {
	Type        string                `json:"type"`                  // "connected", "tool_call", "rag_results", "audio", "text", "transcript", "interrupted", "turn_complete", "error", "pong"
	Data        string                `json:"data,omitempty"`        // Base64 audio PCM (24kHz mono)
	MimeType    string                `json:"mime_type,omitempty"`
	Text        string                `json:"text,omitempty"`
	Query       string                `json:"query,omitempty"`
	ToolName    string                `json:"tool_name,omitempty"`
	CallID      string                `json:"call_id,omitempty"`
	Citations   []MagisteriumCitation `json:"citations,omitempty"`
	Count       int                   `json:"count,omitempty"`
	Error       string                `json:"error,omitempty"`
}

// SafeWebSocket wraps gorilla/websocket to ensure concurrent writes are thread-safe.
type SafeWebSocket struct {
	conn *websocket.Conn
	mu   sync.Mutex
}

func (s *SafeWebSocket) WriteJSON(v interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.conn.WriteJSON(v)
}

func (s *SafeWebSocket) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.conn.Close()
}

// geminiLiveHandler manages real-time bidirectional audio & text with Gemini Live and Magisterium RAG tool calling.
func geminiLiveHandler(w http.ResponseWriter, r *http.Request) {
	// CORS preflight support
	if r.Method == http.MethodOptions {
		allowedOrigin := os.Getenv("CORS_ALLOWED_ORIGIN")
		if allowedOrigin == "" {
			allowedOrigin = "*"
		}
		w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.WriteHeader(http.StatusOK)
		return
	}

	// Validate API token via Authorization header or ?token= query parameter (WebSockets in browser)
	token := r.URL.Query().Get("token")
	if token == "" {
		authHeader := r.Header.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		}
	}

	expectedToken := os.Getenv("API_TOKEN")
	if expectedToken != "" && token != expectedToken {
		log.Printf("[geminiLiveHandler] Unauthorized WebSocket connection attempt (invalid token)")
		http.Error(w, "Unauthorized: Invalid API token", http.StatusUnauthorized)
		return
	}

	settings := GetSettings()
	geminiKey := settings.GeminiAPIKey
	if geminiKey == "" {
		geminiKey = os.Getenv("GEMINI_API_KEY")
	}

	magisteriumKey := settings.MagisteriumAPIKey
	if magisteriumKey == "" {
		magisteriumKey = os.Getenv("MAGISTERIUM_API_KEY")
	}

	if geminiKey == "" {
		http.Error(w, "Gemini API key is not configured in Admin Settings", http.StatusBadRequest)
		return
	}

	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[geminiLiveHandler] Failed to upgrade WebSocket connection: %v", err)
		return
	}

	safeWS := &SafeWebSocket{conn: wsConn}
	defer safeWS.Close()

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	log.Println("[geminiLiveHandler] Client connected to Magisterium Gemini Live Gateway")

	// Determine Live Model (defaults to gemini-3.8-live, independent of text gemini_model)
	liveModel := os.Getenv("GEMINI_LIVE_MODEL")
	if liveModel == "" {
		liveModel = strings.TrimSpace(settings.GeminiLiveModel)
	}
	if liveModel == "" {
		liveModel = "gemini-3.8-live"
	}

	log.Printf("[geminiLiveHandler] Initializing Gemini Live Session with model: %s", liveModel)

	genaiClient, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: geminiKey,
	})
	if err != nil {
		log.Printf("[geminiLiveHandler] Error creating GenAI client: %v", err)
		safeWS.WriteJSON(ServerLiveMessage{
			Type:  "error",
			Error: fmt.Sprintf("Failed to initialize Gemini GenAI client: %v", err),
		})
		return
	}

	systemPrompt := `You are Magisterium AI Live, a faithful, knowledgeable, and reverent Catholic teaching assistant.
Your mission is to assist believers, scholars, and seekers by answering questions strictly grounded in authentic Catholic Sacred Scripture, Sacred Tradition, and Church Magisterium.

CRITICAL INSTRUCTION - CALLING THE MAGISTERIUM RAG AGENT:
Whenever the user asks ANY question regarding Catholic faith, theology, morality, church doctrine, catechism, scriptures, church fathers, liturgy, prayer, or ecumenical councils, YOU MUST ALWAYS CALL the tool "get_magisterium_rag" with a concise, targeted search query.
Do NOT attempt to guess or rely solely on memory when Catholic teaching is involved.

AFTER RECEIVING RAG RESULTS FROM THE AGENT:
1. Synthesize and summarize the retrieved citations, doctrine, and excerpts faithfully.
2. Directly answer the user's question with doctrinal fidelity, clarity, and pastoral warmth.
3. Explicitly reference the cited sources inline (e.g., "According to the Catechism of the Catholic Church paragraph 2559...", "In the encyclical Fides et Ratio, Pope Saint John Paul II explains...").
4. Spoken delivery must be clear, natural, and reverent. Keep spoken responses cohesive and suitably concise for listening.`

	tools := []*genai.Tool{
		{
			FunctionDeclarations: []*genai.FunctionDeclaration{
				{
					Name:        "get_magisterium_rag",
					Description: "Query the official Magisterium AI RAG service to search Catholic Church teachings, Catechism, Papal Encyclicals, Vatican & Ecumenical Council documents, and Sacred Scripture. Call this whenever the user asks any question about Catholic doctrine, faith, morality, prayers, or Church teaching.",
					Parameters: &genai.Schema{
						Type: genai.TypeObject,
						Properties: map[string]*genai.Schema{
							"query": {
								Type:        genai.TypeString,
								Description: "The Catholic doctrinal, theological, or scriptural query to search in the Magisterium database.",
							},
						},
						Required: []string{"query"},
					},
				},
			},
		},
	}

	liveConfig := &genai.LiveConnectConfig{
		ResponseModalities: []genai.Modality{genai.ModalityAudio},
		SpeechConfig: &genai.SpeechConfig{
			VoiceConfig: &genai.VoiceConfig{
				PrebuiltVoiceConfig: &genai.PrebuiltVoiceConfig{
					VoiceName: "Aoede",
				},
			},
		},
		SystemInstruction:        genai.NewContentFromText(systemPrompt, "system"),
		Tools:                    tools,
		OutputAudioTranscription: &genai.AudioTranscriptionConfig{},
		InputAudioTranscription:  &genai.AudioTranscriptionConfig{},
	}

	liveSession, err := genaiClient.Live.Connect(ctx, liveModel, liveConfig)
	if err != nil {
		log.Printf("[geminiLiveHandler] Failed to connect to Gemini Live API: %v", err)
		safeWS.WriteJSON(ServerLiveMessage{
			Type:  "error",
			Error: fmt.Sprintf("Failed to establish Gemini Live session: %v", err),
		})
		return
	}
	defer liveSession.Close()

	log.Println("[geminiLiveHandler] Gemini Live session connected successfully")
	safeWS.WriteJSON(ServerLiveMessage{
		Type: "connected",
		Text: "Magisterium Gemini Live session ready. Speak or send a question.",
	})

	// WaitGroup to manage goroutines
	var wg sync.WaitGroup
	wg.Add(2)

	// Goroutine 1: Read messages from Frontend Client WebSocket and send to Gemini Live
	go func() {
		defer wg.Done()
		defer cancel()

		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			_, rawBytes, err := wsConn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
					log.Printf("[geminiLiveHandler] Client WebSocket closed: %v", err)
				}
				return
			}

			var clientMsg ClientLiveMessage
			if err := json.Unmarshal(rawBytes, &clientMsg); err != nil {
				log.Printf("[geminiLiveHandler] Invalid client message JSON: %v", err)
				continue
			}

			switch clientMsg.Type {
			case "realtime_input":
				if clientMsg.Data != "" {
					audioBytes, err := base64.StdEncoding.DecodeString(clientMsg.Data)
					if err != nil {
						log.Printf("[geminiLiveHandler] Error decoding base64 audio: %v", err)
						continue
					}
					err = liveSession.SendRealtimeInput(genai.LiveSendRealtimeInputParameters{
						Audio: &genai.Blob{
							MIMEType: "audio/pcm;rate=16000",
							Data:     audioBytes,
						},
					})
					if err != nil {
						log.Printf("[geminiLiveHandler] Error sending audio to Gemini Live: %v", err)
					}
				}

			case "client_content":
				if strings.TrimSpace(clientMsg.Text) != "" {
					log.Printf("[geminiLiveHandler] Sending user text prompt to Gemini Live: %q", clientMsg.Text)
					err = liveSession.SendClientContent(genai.LiveSendClientContentParameters{
						Turns: []*genai.Content{
							genai.NewContentFromText(clientMsg.Text, "user"),
						},
						TurnComplete: genai.Ptr(true),
					})
					if err != nil {
						log.Printf("[geminiLiveHandler] Error sending client content to Gemini Live: %v", err)
					}
				}

			case "audio_stream_end":
				log.Println("[geminiLiveHandler] Audio stream end signaled by client")
				err = liveSession.SendRealtimeInput(genai.LiveSendRealtimeInputParameters{
					AudioStreamEnd: true,
				})
				if err != nil {
					log.Printf("[geminiLiveHandler] Error sending audio stream end: %v", err)
				}

			case "ping":
				safeWS.WriteJSON(ServerLiveMessage{Type: "pong"})
			}
		}
	}()

	// Goroutine 2: Receive messages from Gemini Live and stream to Frontend Client
	go func() {
		defer wg.Done()
		defer cancel()

		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			serverMsg, err := liveSession.Receive()
			if err != nil {
				log.Printf("[geminiLiveHandler] Error receiving from Gemini Live session: %v", err)
				safeWS.WriteJSON(ServerLiveMessage{
					Type:  "error",
					Error: fmt.Sprintf("Gemini Live session error: %v", err),
				})
				return
			}

			// Handle Tool Calls (e.g. get_magisterium_rag)
			if serverMsg.ToolCall != nil && len(serverMsg.ToolCall.FunctionCalls) > 0 {
				var functionResponses []*genai.FunctionResponse

				for _, call := range serverMsg.ToolCall.FunctionCalls {
					log.Printf("[geminiLiveHandler] Gemini Live requested tool call: %s (ID: %s)", call.Name, call.ID)

					if call.Name == "get_magisterium_rag" {
						queryVal, _ := call.Args["query"].(string)
						if queryVal == "" {
							queryVal = "Catholic Church teaching"
						}

						log.Printf("[geminiLiveHandler] Agent executing Magisterium RAG for: %q", queryVal)

						// Inform client that tool is being invoked
						safeWS.WriteJSON(ServerLiveMessage{
							Type:     "tool_call",
							ToolName: "get_magisterium_rag",
							Query:    queryVal,
							CallID:   call.ID,
						})

						// Execute RAG Agent
						citations, _, ragErr := ExecuteMagisteriumRAG(ctx, queryVal, 5, magisteriumKey)
						if ragErr != nil {
							log.Printf("[geminiLiveHandler] Magisterium RAG Agent error: %v", ragErr)
						}

						// Forward retrieved citations to frontend client for real-time display
						safeWS.WriteJSON(ServerLiveMessage{
							Type:      "rag_results",
							Query:     queryVal,
							Citations: citations,
							Count:     len(citations),
						})

						formattedContext := FormatCitationsAsContext(citations)
						if ragErr != nil && len(citations) == 0 {
							formattedContext = fmt.Sprintf("Magisterium search error: %v. Please advise the user accordingly.", ragErr)
						}

						functionResponses = append(functionResponses, &genai.FunctionResponse{
							ID:   call.ID,
							Name: call.Name,
							Response: map[string]any{
								"query":           queryVal,
								"total_citations": len(citations),
								"citations":       citations,
								"context":         formattedContext,
							},
						})
					} else {
						functionResponses = append(functionResponses, &genai.FunctionResponse{
							ID:   call.ID,
							Name: call.Name,
							Response: map[string]any{
								"error": fmt.Sprintf("Unknown tool %s", call.Name),
							},
						})
					}
				}

				// Return tool execution outputs back to Gemini Live
				if len(functionResponses) > 0 {
					err = liveSession.SendToolResponse(genai.LiveSendToolResponseParameters{
						FunctionResponses: functionResponses,
					})
					if err != nil {
						log.Printf("[geminiLiveHandler] Error sending tool response to Gemini Live: %v", err)
					}
				}
			}

			// Handle Server Content (Spoken Audio, Text, Transcripts, Interrupted)
			if serverMsg.ServerContent != nil {
				sc := serverMsg.ServerContent

				if sc.Interrupted {
					log.Println("[geminiLiveHandler] Gemini Live playback interrupted by user speech")
					safeWS.WriteJSON(ServerLiveMessage{
						Type: "interrupted",
					})
				}

				if sc.ModelTurn != nil {
					for _, part := range sc.ModelTurn.Parts {
						// Stream real-time spoken audio chunks
						if part.InlineData != nil && len(part.InlineData.Data) > 0 {
							safeWS.WriteJSON(ServerLiveMessage{
								Type:     "audio",
								Data:     base64.StdEncoding.EncodeToString(part.InlineData.Data),
								MimeType: part.InlineData.MIMEType,
							})
						}

						// Stream text parts if available
						if part.Text != "" {
							safeWS.WriteJSON(ServerLiveMessage{
								Type: "text",
								Text: part.Text,
							})
						}
					}
				}

				// Stream transcript text
				if sc.OutputTranscription != nil && sc.OutputTranscription.Text != "" {
					safeWS.WriteJSON(ServerLiveMessage{
						Type: "transcript",
						Text: sc.OutputTranscription.Text,
					})
				}

				if sc.TurnComplete {
					log.Println("[geminiLiveHandler] Gemini Live turn complete")
					safeWS.WriteJSON(ServerLiveMessage{
						Type: "turn_complete",
					})
				}
			}

			if serverMsg.GoAway != nil {
				log.Println("[geminiLiveHandler] Gemini Live server sent GoAway")
				safeWS.WriteJSON(ServerLiveMessage{
					Type: "error",
					Text: "Gemini Live session disconnected by server.",
				})
				return
			}
		}
	}()

	// Keep alive ping loop
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				safeWS.WriteJSON(ServerLiveMessage{Type: "ping"})
			}
		}
	}()

	wg.Wait()
	log.Println("[geminiLiveHandler] Gemini Live session terminated gracefully")
}
