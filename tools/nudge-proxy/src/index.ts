import { Hono } from 'hono'
import { cors } from 'hono/cors'
import { serve } from '@hono/node-server'
import OpenAI from 'openai'

const app = new Hono()

app.use('/*', cors())

app.get('/', (c) => c.text('LunaFox Nudge Proxy is running! 🦊'))

app.post('/generate', async (c) => {
  try {
    const { context } = await c.req.json()
    const apiKey = process.env.AI_API_KEY
    const baseURL = process.env.AI_BASE_URL || 'https://api.deepseek.com/v1'

    if (!apiKey) {
      return c.json({ error: 'Server AI_API_KEY not configured' }, 500)
    }

    const openai = new OpenAI({ apiKey, baseURL })

    const prompt = `
You are LunaFox (月狐), a witty, humorous, hacker-culture-savvy virtual assistant for cybersecurity professionals.
Current context:
- Hour: ${context.hour}
- Day: ${context.day} (0=Sun, 5=Fri)
- Event: ${context.event || 'idle'}

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
`

    const completion = await openai.chat.completions.create({
      messages: [{ role: 'system', content: prompt }],
      model: 'deepseek-chat', // Default to deepseek, configurable
      response_format: { type: 'json_object' },
    })

    const content = completion.choices[0].message.content
    if (!content) throw new Error('No content from AI')

    const data = JSON.parse(content)
    return c.json(data)
  } catch (err) {
    console.error(err)
    return c.json({ error: 'Failed to generate nudge' }, 500)
  }
})

const port = Number(process.env.PORT) || 3000
console.log(`Server is running on port ${port}`)

serve({
  fetch: app.fetch,
  port
})
