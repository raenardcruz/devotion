import { ref } from 'vue';

export interface MagisteriumLiveCitation {
  title?: string;
  document_title?: string;
  author?: string;
  ref?: string;
  score?: number;
  url?: string;
  text?: string;
}

export interface LiveTurnResult {
  userQuery?: string;
  response: string;
  citations: MagisteriumLiveCitation[];
}

export function useGeminiLive(onTurnComplete?: (result: LiveTurnResult) => void) {
  const isLiveConnected = ref(false);
  const isConnecting = ref(false);
  const isMuted = ref(false);
  const isGeminiSpeaking = ref(false);
  const isUserSpeaking = ref(false);
  const isQueryingRAG = ref(false);
  const currentRagQuery = ref<string>('');
  const currentCitations = ref<MagisteriumLiveCitation[]>([]);
  const liveTranscript = ref<string>('');
  const userLiveTranscript = ref<string>('');
  const liveStatus = ref<string>('Idle');
  const liveError = ref<string | null>(null);
  const audioLevel = ref<number>(0);

  let ws: WebSocket | null = null;
  let audioContext: AudioContext | null = null;
  let micMediaStream: MediaStream | null = null;
  let micSourceNode: MediaStreamAudioSourceNode | null = null;
  let processorNode: ScriptProcessorNode | null = null;
  let playbackAudioContext: AudioContext | null = null;
  let nextAudioStartTime = 0;
  const activeAudioSources: AudioBufferSourceNode[] = [];

  // Helper to convert base64 to Uint8Array
  const base64ToUint8Array = (base64: string): Uint8Array => {
    const binary = atob(base64);
    const len = binary.length;
    const bytes = new Uint8Array(len);
    for (let i = 0; i < len; i++) {
      bytes[i] = binary.charCodeAt(i);
    }
    return bytes;
  };

  // Helper to convert Float32Array to 16-bit PCM Int16Array
  const floatTo16BitPCM = (input: Float32Array): Int16Array => {
    const output = new Int16Array(input.length);
    for (let i = 0; i < input.length; i++) {
      const val = input[i] ?? 0;
      const s = Math.max(-1, Math.min(1, val));
      output[i] = s < 0 ? s * 0x8000 : s * 0x7FFF;
    }
    return output;
  };

  // Helper to encode ArrayBuffer to base64
  const arrayBufferToBase64 = (buffer: ArrayBufferLike): string => {
    let binary = '';
    const bytes = new Uint8Array(buffer as ArrayBuffer);
    const len = bytes.byteLength;
    for (let i = 0; i < len; i++) {
      const byte = bytes[i] ?? 0;
      binary += String.fromCharCode(byte);
    }
    return btoa(binary);
  };

  // Resample audio buffer to 16kHz
  const downsampleTo16k = (input: Float32Array, inputRate: number): Float32Array => {
    if (inputRate === 16000) return input;
    const ratio = inputRate / 16000;
    const newLength = Math.round(input.length / ratio);
    const result = new Float32Array(newLength);
    let offsetResult = 0;
    let offsetInput = 0;
    while (offsetResult < result.length) {
      const nextOffsetInput = Math.round((offsetResult + 1) * ratio);
      let accum = 0;
      let count = 0;
      for (let i = offsetInput; i < nextOffsetInput && i < input.length; i++) {
        accum += input[i] ?? 0;
        count++;
      }
      result[offsetResult] = count > 0 ? accum / count : 0;
      offsetResult++;
      offsetInput = nextOffsetInput;
    }
    return result;
  };

  // Playback PCM 24kHz audio chunk from Gemini
  const playAudioChunk = (base64Pcm: string) => {
    if (!playbackAudioContext) {
      playbackAudioContext = new (window.AudioContext || (window as any).webkitAudioContext)({ sampleRate: 24000 });
    }

    if (playbackAudioContext.state === 'suspended') {
      playbackAudioContext.resume();
    }

    const uint8 = base64ToUint8Array(base64Pcm);
    const int16 = new Int16Array(uint8.buffer, uint8.byteOffset, uint8.byteLength / 2);
    const float32 = new Float32Array(int16.length);

    for (let i = 0; i < int16.length; i++) {
      const val = int16[i] ?? 0;
      float32[i] = val / 32768.0;
    }

    const audioBuffer = playbackAudioContext.createBuffer(1, float32.length, 24000);
    audioBuffer.getChannelData(0).set(float32);

    const source = playbackAudioContext.createBufferSource();
    source.buffer = audioBuffer;
    source.connect(playbackAudioContext.destination);

    const now = playbackAudioContext.currentTime;
    const startTime = Math.max(now, nextAudioStartTime);
    source.start(startTime);
    nextAudioStartTime = startTime + audioBuffer.duration;

    isGeminiSpeaking.value = true;
    liveStatus.value = 'Gemini Live speaking...';

    activeAudioSources.push(source);
    source.onended = () => {
      const idx = activeAudioSources.indexOf(source);
      if (idx !== -1) {
        activeAudioSources.splice(idx, 1);
      }
      if (activeAudioSources.length === 0 && (!playbackAudioContext || playbackAudioContext.currentTime >= nextAudioStartTime - 0.05)) {
        isGeminiSpeaking.value = false;
        if (isLiveConnected.value && !isQueryingRAG.value) {
          liveStatus.value = 'Listening...';
        }
      }
    };
  };

  // Interrupt / stop all active playback
  const interruptPlayback = () => {
    activeAudioSources.forEach(s => {
      try {
        s.stop();
        s.disconnect();
      } catch {}
    });
    activeAudioSources.length = 0;
    nextAudioStartTime = 0;
    isGeminiSpeaking.value = false;
  };

  // Start capturing audio from microphone
  const startMicrophone = async () => {
    try {
      micMediaStream = await navigator.mediaDevices.getUserMedia({
        audio: {
          channelCount: 1,
          echoCancellation: true,
          noiseSuppression: true,
          autoGainControl: true,
        },
      });

      audioContext = new (window.AudioContext || (window as any).webkitAudioContext)();
      if (audioContext.state === 'suspended') {
        await audioContext.resume();
      }

      micSourceNode = audioContext.createMediaStreamSource(micMediaStream);
      processorNode = audioContext.createScriptProcessor(4096, 1, 1);

      processorNode.onaudioprocess = (e) => {
        if (!isLiveConnected.value || isMuted.value || !ws || ws.readyState !== WebSocket.OPEN) {
          audioLevel.value = 0;
          return;
        }

        const inputData = e.inputBuffer.getChannelData(0);

        // Compute volume level (RMS)
        let sum = 0;
        for (let i = 0; i < inputData.length; i++) {
          const sample = inputData[i] ?? 0;
          sum += sample * sample;
        }
        const rms = Math.sqrt(sum / inputData.length);
        audioLevel.value = Math.min(1, rms * 5);
        isUserSpeaking.value = rms > 0.02;

        // Downsample to 16kHz
        const downsampled = downsampleTo16k(inputData, audioContext!.sampleRate);
        const pcm16 = floatTo16BitPCM(downsampled);
        const base64Data = arrayBufferToBase64(pcm16.buffer);

        ws.send(JSON.stringify({
          type: 'realtime_input',
          data: base64Data,
        }));
      };

      micSourceNode.connect(processorNode);
      processorNode.connect(audioContext.destination);
    } catch (err: any) {
      console.error('[useGeminiLive] Microphone access error:', err);
      liveError.value = 'Failed to access microphone. Please grant audio permissions in your browser.';
    }
  };

  const stopMicrophone = () => {
    if (processorNode) {
      processorNode.disconnect();
      processorNode = null;
    }
    if (micSourceNode) {
      micSourceNode.disconnect();
      micSourceNode = null;
    }
    if (micMediaStream) {
      micMediaStream.getTracks().forEach(t => t.stop());
      micMediaStream = null;
    }
    if (audioContext && audioContext.state !== 'closed') {
      audioContext.close();
      audioContext = null;
    }
    audioLevel.value = 0;
    isUserSpeaking.value = false;
  };

  const toggleMute = () => {
    isMuted.value = !isMuted.value;
    if (isMuted.value) {
      audioLevel.value = 0;
      isUserSpeaking.value = false;
      if (ws && ws.readyState === WebSocket.OPEN) {
        ws.send(JSON.stringify({ type: 'audio_stream_end' }));
      }
    }
  };

  // Connect to Gemini Live Gateway over WebSocket
  const startLiveSession = async () => {
    if (isLiveConnected.value || isConnecting.value) return;

    liveError.value = null;
    isConnecting.value = true;
    liveStatus.value = 'Connecting to Gemini Live...';
    currentCitations.value = [];
    liveTranscript.value = '';
    userLiveTranscript.value = '';

    try {
      let baseUrl = import.meta.env.VITE_API_BASE_URL || 'https://devotionapi.raenardcruz.com';
      if (!baseUrl || (!baseUrl.startsWith('http://') && !baseUrl.startsWith('https://') && !baseUrl.startsWith('/'))) {
        baseUrl = 'https://devotionapi.raenardcruz.com';
      }

      let wsUrl = baseUrl.replace(/^http:\/\//, 'ws://').replace(/^https:\/\//, 'wss://');
      if (wsUrl.endsWith('/')) {
        wsUrl = wsUrl.slice(0, -1);
      }
      wsUrl += '/magisterium/live';

      const apiToken = import.meta.env.VITE_API_TOKEN;
      if (apiToken) {
        wsUrl += `?token=${encodeURIComponent(apiToken)}`;
      }

      console.log('[useGeminiLive] Connecting to Live WebSocket:', wsUrl.replace(/\?token=.*/, '?token=***'));
      ws = new WebSocket(wsUrl);

      ws.onopen = async () => {
        console.log('[useGeminiLive] WebSocket connected');
        isLiveConnected.value = true;
        isConnecting.value = false;
        liveStatus.value = 'Connected. Starting microphone...';
        await startMicrophone();
        liveStatus.value = 'Listening...';
      };

      ws.onmessage = (event) => {
        try {
          const msg = JSON.parse(event.data);

          switch (msg.type) {
            case 'connected':
              liveStatus.value = 'Listening...';
              break;

            case 'tool_call':
              isQueryingRAG.value = true;
              currentRagQuery.value = msg.query || '';
              liveStatus.value = `Agent querying Magisterium RAG for: "${msg.query}"...`;
              break;

            case 'rag_results':
              isQueryingRAG.value = false;
              if (msg.citations && msg.citations.length > 0) {
                currentCitations.value = msg.citations;
                liveStatus.value = `Retrieved ${msg.citations.length} Magisterial citations. Summarizing...`;
              }
              break;

            case 'audio':
              if (msg.data) {
                playAudioChunk(msg.data);
              }
              break;

            case 'text':
              if (msg.text) {
                liveTranscript.value += msg.text;
              }
              break;

            case 'transcript':
              if (msg.text) {
                liveTranscript.value += msg.text;
              }
              break;

            case 'interrupted':
              console.log('[useGeminiLive] Interrupted by user');
              interruptPlayback();
              liveStatus.value = 'Listening...';
              break;

            case 'turn_complete':
              console.log('[useGeminiLive] Turn complete');
              isGeminiSpeaking.value = false;
              liveStatus.value = 'Listening...';

              if (liveTranscript.value.trim() && onTurnComplete) {
                onTurnComplete({
                  userQuery: userLiveTranscript.value.trim() || currentRagQuery.value || 'Voice Question',
                  response: liveTranscript.value.trim(),
                  citations: [...currentCitations.value],
                });
              }

              // Reset per-turn transcript buffers
              liveTranscript.value = '';
              userLiveTranscript.value = '';
              break;

            case 'error':
              console.error('[useGeminiLive] Server reported error:', msg.error);
              liveError.value = msg.error || 'Gemini Live encountered an error.';
              liveStatus.value = 'Error';
              break;

            case 'pong':
              break;
          }
        } catch (e) {
          console.error('[useGeminiLive] Error parsing incoming WebSocket message:', e);
        }
      };

      ws.onerror = (e) => {
        console.error('[useGeminiLive] WebSocket error:', e);
        liveError.value = 'Failed to connect to Magisterium Gemini Live Gateway.';
        liveStatus.value = 'Connection error';
      };

      ws.onclose = () => {
        console.log('[useGeminiLive] WebSocket closed');
        cleanup();
      };
    } catch (err: any) {
      console.error('[useGeminiLive] Live session error:', err);
      liveError.value = err.message || 'Failed to start Gemini Live session.';
      cleanup();
    }
  };

  // Send a text message turn to Gemini Live
  const sendTextMessage = (text: string) => {
    if (!ws || ws.readyState !== WebSocket.OPEN) return;
    const trimmed = text.trim();
    if (!trimmed) return;

    interruptPlayback();
    userLiveTranscript.value = trimmed;
    liveTranscript.value = '';
    currentCitations.value = [];
    liveStatus.value = 'Processing question...';

    ws.send(JSON.stringify({
      type: 'client_content',
      text: trimmed,
    }));
  };

  // Stop / Disconnect Live session
  const stopLiveSession = () => {
    if (ws) {
      try {
        ws.close();
      } catch {}
      ws = null;
    }
    cleanup();
  };

  const cleanup = () => {
    stopMicrophone();
    interruptPlayback();
    if (playbackAudioContext && playbackAudioContext.state !== 'closed') {
      try {
        playbackAudioContext.close();
      } catch {}
      playbackAudioContext = null;
    }
    isLiveConnected.value = false;
    isConnecting.value = false;
    isGeminiSpeaking.value = false;
    isQueryingRAG.value = false;
    liveStatus.value = 'Disconnected';
  };

  return {
    isLiveConnected,
    isConnecting,
    isMuted,
    isGeminiSpeaking,
    isUserSpeaking,
    isQueryingRAG,
    currentRagQuery,
    currentCitations,
    liveTranscript,
    userLiveTranscript,
    liveStatus,
    liveError,
    audioLevel,
    startLiveSession,
    stopLiveSession,
    toggleMute,
    sendTextMessage,
    interruptPlayback,
  };
}
