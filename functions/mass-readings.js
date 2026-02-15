import { GoogleGenerativeAI } from "@google/generative-ai";

export async function onRequest(context) {
    const apiKey = context.env.GEMINI_API_KEY;

    if (!apiKey) {
        return new Response("Missing GEMINI_API_KEY", { status: 500 });
    }

    const genAI = new GoogleGenerativeAI(apiKey);
    const model = genAI.getGenerativeModel({ model: "gemini-2.5-flash" });

    const today = new Date().toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'long',
        day: 'numeric'
    });

    const prompt = `Give me the roman readings for ${today} in nrsvce tranlation and the historical/literal context. Make sure to get the Year A, B, or C depending on the year to get the correct mass reading. Set second_reading as null if not applicable. Output only the json format: { first_reading: { verse: '', text: '', context: '' }, responsorial_psalm: { verse: '', text: '' }, second_reading: { verse: '', text: '', context: '' }, gospel: { verse: '', text: '', context: '' } }`;

    console.log("Generated Prompt:", prompt);

    try {
        const result = await model.generateContent(prompt);
        const response = await result.response;
        const text = response.text();

        // Clean up the response to ensure it's valid JSON
        const jsonString = text.replace(/```json/g, '').replace(/```/g, '').trim();

        return new Response(jsonString, {
            headers: {
                "Content-Type": "application/json",
                "Access-Control-Allow-Origin": "*", // Allow CORS for local dev
            },
        });
    } catch (error) {
        return new Response(JSON.stringify({ error: error.message }), {
            status: 500,
            headers: { "Content-Type": "application/json" },
        });
    }
}
