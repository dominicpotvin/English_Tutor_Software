/**
 * MCP connector for the English Tutor platform.
 *
 * Exposes the platform's curriculum and progress over the Model Context
 * Protocol so an assistant can review and co-author lessons, topics,
 * exercises, quizzes and vocabulary alongside the learner.
 *
 * It is a thin client over the backend REST API; all business rules and
 * validation live in the Go backend.
 */
import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js";
import { StdioServerTransport } from "@modelcontextprotocol/sdk/server/stdio.js";
import { z } from "zod";

const API_URL = process.env.TUTOR_API_URL ?? "http://localhost:8096/api";

/** apiCall performs a JSON request against the backend and returns the parsed body. */
async function apiCall(method: string, path: string, body?: unknown): Promise<unknown> {
  const res = await fetch(API_URL + path, {
    method,
    headers: { "Content-Type": "application/json" },
    body: body === undefined ? undefined : JSON.stringify(body),
  });
  const raw = await res.text();
  if (!res.ok) {
    throw new Error(`${method} ${path} -> HTTP ${res.status}: ${raw || res.statusText}`);
  }
  return raw ? JSON.parse(raw) : null;
}

/** result wraps a value as an MCP text result. */
function result(data: unknown) {
  return { content: [{ type: "text" as const, text: JSON.stringify(data, null, 2) }] };
}

const server = new McpServer({ name: "english-tutor", version: "1.0.0" });

// --- Reading the curriculum ----------------------------------------------

server.registerTool(
  "list_curriculum",
  {
    title: "List curriculum",
    description:
      "List every level with its lessons (id, number, title, topic and exercise counts). " +
      "Use this first to find the id of the lesson or topic you want to work on.",
    inputSchema: {},
  },
  async () => {
    const levels = (await apiCall("GET", "/levels")) as Array<Record<string, unknown>>;
    const tree = [];
    for (const level of levels) {
      const detail = (await apiCall("GET", `/levels/${level.id}`)) as {
        lessons: Array<Record<string, unknown>>;
      };
      tree.push({
        levelId: level.id,
        code: level.code,
        name: level.name,
        lessons: detail.lessons.map((l) => ({
          lessonId: l.id,
          number: l.number,
          title: l.title,
          topicCount: l.topicCount,
          exerciseCount: l.exerciseCount,
        })),
      });
    }
    return result(tree);
  },
);

server.registerTool(
  "get_lesson",
  {
    title: "Get lesson",
    description: "Return a lesson with its level and all of its topics, including the Markdown explanations.",
    inputSchema: { lessonId: z.number().int().describe("Lesson id") },
  },
  async ({ lessonId }) => result(await apiCall("GET", `/lessons/${lessonId}`)),
);

server.registerTool(
  "get_topic",
  {
    title: "Get topic",
    description: "Return a topic with all of its exercises, including answers and explanations.",
    inputSchema: { topicId: z.number().int().describe("Topic id") },
  },
  async ({ topicId }) => result(await apiCall("GET", `/topics/${topicId}/exercises`)),
);

// --- Lessons --------------------------------------------------------------

server.registerTool(
  "create_lesson",
  {
    title: "Create lesson",
    description: "Add a lesson to a level.",
    inputSchema: {
      levelId: z.number().int(),
      number: z.number().int().describe("Lesson number within the level"),
      title: z.string(),
      summary: z.string().optional(),
      position: z.number().int().optional(),
    },
  },
  async (args) =>
    result(
      await apiCall("POST", "/lessons", {
        levelId: args.levelId,
        number: args.number,
        title: args.title,
        summary: args.summary ?? "",
        position: args.position ?? args.number,
      }),
    ),
);

server.registerTool(
  "update_lesson",
  {
    title: "Update lesson",
    description: "Update a lesson. Only the fields you provide change.",
    inputSchema: {
      lessonId: z.number().int(),
      number: z.number().int().optional(),
      title: z.string().optional(),
      summary: z.string().optional(),
      position: z.number().int().optional(),
    },
  },
  async (args) => {
    const data = (await apiCall("GET", `/lessons/${args.lessonId}`)) as {
      lesson: Record<string, unknown>;
    };
    const cur = data.lesson;
    return result(
      await apiCall("PUT", `/lessons/${args.lessonId}`, {
        levelId: cur.levelId,
        number: args.number ?? cur.number,
        title: args.title ?? cur.title,
        summary: args.summary ?? cur.summary,
        position: args.position ?? cur.position,
      }),
    );
  },
);

// --- Topics ---------------------------------------------------------------

server.registerTool(
  "create_topic",
  {
    title: "Create topic",
    description:
      "Add a teaching topic to a lesson. The explanation is Markdown and may use headings, " +
      "bold text, lists and tables.",
    inputSchema: {
      lessonId: z.number().int(),
      title: z.string(),
      explanation: z.string().describe("Markdown teaching content"),
      position: z.number().int().optional(),
    },
  },
  async (args) =>
    result(
      await apiCall("POST", "/topics", {
        lessonId: args.lessonId,
        title: args.title,
        explanation: args.explanation,
        position: args.position ?? 0,
      }),
    ),
);

server.registerTool(
  "update_topic",
  {
    title: "Update topic",
    description: "Update a topic. Only the fields you provide change.",
    inputSchema: {
      topicId: z.number().int(),
      title: z.string().optional(),
      explanation: z.string().optional(),
      position: z.number().int().optional(),
    },
  },
  async (args) => {
    const cur = (await apiCall("GET", `/topics/${args.topicId}`)) as Record<string, unknown>;
    return result(
      await apiCall("PUT", `/topics/${args.topicId}`, {
        lessonId: cur.lessonId,
        title: args.title ?? cur.title,
        explanation: args.explanation ?? cur.explanation,
        position: args.position ?? cur.position,
      }),
    );
  },
);

// --- Exercises ------------------------------------------------------------

const exerciseKind = z.enum(["mcq", "fill_blank", "true_false"]);

server.registerTool(
  "create_exercise",
  {
    title: "Create exercise",
    description:
      "Add an exercise to a topic (practice) or a quiz (assessment) - provide exactly one of " +
      "topicId or quizId. For mcq the answer must equal one choice exactly; for fill_blank leave " +
      "choices empty and use ___ in the prompt to mark the blank; for true_false use the choices " +
      "True and False.",
    inputSchema: {
      topicId: z.number().int().optional(),
      quizId: z.number().int().optional(),
      kind: exerciseKind,
      prompt: z.string(),
      choices: z.array(z.string()).optional(),
      answer: z.string(),
      explanation: z.string().optional(),
      position: z.number().int().optional(),
    },
  },
  async (args) =>
    result(
      await apiCall("POST", "/exercises", {
        topicId: args.topicId ?? null,
        quizId: args.quizId ?? null,
        kind: args.kind,
        prompt: args.prompt,
        choices: args.choices ?? [],
        answer: args.answer,
        explanation: args.explanation ?? "",
        position: args.position ?? 0,
      }),
    ),
);

server.registerTool(
  "update_exercise",
  {
    title: "Update exercise",
    description: "Update an exercise. Only the fields you provide change.",
    inputSchema: {
      exerciseId: z.number().int(),
      kind: exerciseKind.optional(),
      prompt: z.string().optional(),
      choices: z.array(z.string()).optional(),
      answer: z.string().optional(),
      explanation: z.string().optional(),
      position: z.number().int().optional(),
    },
  },
  async (args) => {
    const cur = (await apiCall("GET", `/exercises/${args.exerciseId}`)) as Record<string, unknown>;
    return result(
      await apiCall("PUT", `/exercises/${args.exerciseId}`, {
        topicId: cur.topicId,
        quizId: cur.quizId,
        kind: args.kind ?? cur.kind,
        prompt: args.prompt ?? cur.prompt,
        choices: args.choices ?? cur.choices,
        answer: args.answer ?? cur.answer,
        explanation: args.explanation ?? cur.explanation,
        position: args.position ?? cur.position,
      }),
    );
  },
);

server.registerTool(
  "delete_exercise",
  {
    title: "Delete exercise",
    description: "Permanently delete an exercise.",
    inputSchema: { exerciseId: z.number().int() },
  },
  async ({ exerciseId }) => {
    await apiCall("DELETE", `/exercises/${exerciseId}`);
    return result({ deleted: exerciseId });
  },
);

// --- Quizzes --------------------------------------------------------------

server.registerTool(
  "list_quizzes",
  {
    title: "List quizzes",
    description: "List every quiz with its question count.",
    inputSchema: {},
  },
  async () => result(await apiCall("GET", "/quizzes")),
);

server.registerTool(
  "get_quiz",
  {
    title: "Get quiz",
    description: "Return a quiz with all of its questions.",
    inputSchema: { quizId: z.number().int() },
  },
  async ({ quizId }) => result(await apiCall("GET", `/quizzes/${quizId}`)),
);

server.registerTool(
  "create_quiz",
  {
    title: "Create quiz",
    description: "Create a quiz. Add questions afterwards with create_exercise using the new quizId.",
    inputSchema: {
      levelId: z.number().int().optional().describe("Level the quiz belongs to"),
      title: z.string(),
      description: z.string().optional(),
      position: z.number().int().optional(),
    },
  },
  async (args) =>
    result(
      await apiCall("POST", "/quizzes", {
        levelId: args.levelId ?? null,
        title: args.title,
        description: args.description ?? "",
        position: args.position ?? 0,
      }),
    ),
);

// --- Vocabulary -----------------------------------------------------------

server.registerTool(
  "list_vocabulary",
  {
    title: "List vocabulary",
    description: "List every vocabulary entry.",
    inputSchema: {},
  },
  async () => result(await apiCall("GET", "/vocabulary")),
);

server.registerTool(
  "add_vocabulary",
  {
    title: "Add vocabulary",
    description: "Add a vocabulary term.",
    inputSchema: {
      levelId: z.number().int().optional(),
      lessonId: z.number().int().optional(),
      category: z.string(),
      term: z.string(),
      definition: z.string().optional(),
      example: z.string().optional(),
      position: z.number().int().optional(),
    },
  },
  async (args) =>
    result(
      await apiCall("POST", "/vocabulary", {
        levelId: args.levelId ?? null,
        lessonId: args.lessonId ?? null,
        category: args.category,
        term: args.term,
        definition: args.definition ?? "",
        example: args.example ?? "",
        position: args.position ?? 0,
      }),
    ),
);

// --- Progress -------------------------------------------------------------

server.registerTool(
  "get_progress",
  {
    title: "Get learner progress",
    description: "Return the learner's progress: mastered exercises, attempts and accuracy per level.",
    inputSchema: {},
  },
  async () => result(await apiCall("GET", "/progress")),
);

// --- Start ----------------------------------------------------------------

await server.connect(new StdioServerTransport());
