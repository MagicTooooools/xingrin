const OPENAI_API_KEY = Deno.env.get("AI_API_KEY");
const OPENAI_BASE_URL = Deno.env.get("AI_BASE_URL") || "https://api.deepseek.com/v1";

const CORS_HEADERS = {
  "Access-Control-Allow-Origin": "*",
  "Access-Control-Allow-Methods": "POST, OPTIONS",
  "Access-Control-Allow-Headers": "Content-Type",
};

// Use modern Deno.serve (no imports required)
Deno.serve(async (req) => {
  // Handle CORS preflight
  if (req.method === "OPTIONS") {
    return new Response(null, { headers: CORS_HEADERS });
  }

  if (req.method !== "POST") {
    return new Response("Only POST is allowed", { status: 405, headers: CORS_HEADERS });
  }

  try {
    if (!OPENAI_API_KEY) {
      return new Response(JSON.stringify({ error: "Server AI_API_KEY not configured" }), {
        status: 500,
        headers: { "Content-Type": "application/json", ...CORS_HEADERS },
      });
    }

    const { context } = await req.json();

    const prompt = `
You are LunaFox (月狐), a witty, humorous, hacker-culture-savvy virtual assistant for cybersecurity professionals.
Current context:
- Hour: ${context?.hour}
- Day: ${context?.day} (0=Sun, 5=Fri)
- Event: ${context?.event || "idle"}

Task: Generate a single JSON object for a toast notification.
Format:
{
  "title": "Short title (<10 chars)",
  "description": "One sentence message (<30 chars), hacker style, maybe funny or warm.",
  "icon": "A single emoji representing the mood",
  "primaryAction": { "label": "Button Text" },
  "secondaryAction": { "label": "Cancel Text" }
}

Do NOT output markdown. Output ONLY the JSON string.
`;

    // Fetch OpenAI/DeepSeek API
    const aiResponse = await fetch(`${OPENAI_BASE_URL}/chat/completions`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${OPENAI_API_KEY}`,
      },
      body: JSON.stringify({
        model: "deepseek-chat", // Or user configured model
        messages: [{ role: "system", content: prompt }],
        response_format: { type: "json_object" },
      }),
    });

    if (!aiResponse.ok) {
      const err = await aiResponse.text();
      throw new Error(`AI API Error: ${err}`);
    }

    const aiData = await aiResponse.json();
    const content = aiData.choices[0].message.content;
    const jsonContent = JSON.parse(content);

    return new Response(JSON.stringify(jsonContent), {
      headers: { "Content-Type": "application/json", ...CORS_HEADERS },
    });
  } catch (err: any) {
    console.error(err);
    return new Response(JSON.stringify({ error: "Failed to generate nudge", details: err.message }), {
      status: 500,
      headers: { "Content-Type": "application/json", ...CORS_HEADERS },
    });
  }
});
