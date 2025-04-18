// src/api.ts
export interface Task {
  id: number;
  title: string;
  description: string;
  priority: number;
  status: number;
  creation_ts: string;
  user_id: number;
}

/**
 * Проверяем, есть ли действующая сессия (кука).
 * Например, пробуем GET /get_tasks,
 * если 200 -> return true,
 * если 401 -> return false,
 * иначе ошибка.
 */
export async function checkAuth(): Promise<boolean> {
  const res = await fetch('/get_tasks', {
    method: 'GET',
    credentials: 'include',
  });
  if (res.ok) {
    return true;
  } else if (res.status === 401) {
    return false;
  } else {
    throw new Error(`checkAuth failed: ${res.statusText}`);
  }
}

export async function signIn(username: string, password: string): Promise<void> {
  const res = await fetch('/sign_in', {
    method: 'POST',
    credentials: 'include',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, password }),
  });
  if (!res.ok) {
    throw new Error(`Sign in failed: ${res.statusText}`);
  }
}

export async function signUp(username: string, password: string, email: string): Promise<void> {
  const res = await fetch('/sign_up', {
    method: 'POST',
    credentials: 'include',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, password, email }),
  });
  if (!res.ok) {
    throw new Error(`Sign up failed: ${res.statusText}`);
  }
}

export async function logout(): Promise<void> {
  const res = await fetch('/logout', {
    method: 'POST',
    credentials: 'include',
  });
  if (!res.ok) {
    throw new Error(`Logout failed: ${res.statusText}`);
  }
}

export async function getTasks(): Promise<Task[]> {
  const res = await fetch('/get_tasks', {
    method: 'GET',
    credentials: 'include',
  });
  if (!res.ok) {
    throw new Error(`getTasks failed: ${res.statusText}`);
  }
  const data = await res.json();
  return data.tasks;
}

export async function createTask(title: string, description: string): Promise<Task> {
  const body = {
    task: { title, description },
  };
  const res = await fetch('/create_task', {
    method: 'POST',
    credentials: 'include',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  });
  if (!res.ok) {
    throw new Error(`createTask failed: ${res.statusText}`);
  }
  const data = await res.json();
  // data.task — это новый Task
  return data.task;
}


// src/api.ts

// Описание тела запроса
interface UpdatePriorityRequest {
  target_task: Task;
  prev_task_priority?: number;
  next_task_priority?: number;
}

// Описание тела ответа
interface UpdatePriorityResponse {
  task: Task;
}

/**
 * Функция для обновления приоритета задачи
 * (вызывается при Drag & Drop)
 */
export async function updatePriority(
  targetTask: Task,
  prevPriority?: number,
  nextPriority?: number
): Promise<Task> {
  const reqBody: UpdatePriorityRequest = {
    target_task: targetTask,
    prev_task_priority: prevPriority,
    next_task_priority: nextPriority,
  };

  const res = await fetch('/update_priority', {
    method: 'POST',
    credentials: 'include', // чтобы сервер получал и отсылал cookie
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(reqBody),
  });

  if (!res.ok) {
    throw new Error(`updatePriority failed: ${res.statusText}`);
  }

  // По OpenAPI, эндпоинт возвращает handlers.UpdatePriorityResponse
  // { "task": {...} }
  const data = (await res.json()) as UpdatePriorityResponse;
  return data.task;
}

// src/api.ts
export async function updateTask(updatedTask: Task): Promise<Task> {
  const reqBody = { task: updatedTask };
  const res = await fetch('/update_task', {
    method: 'POST',
    credentials: 'include',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(reqBody),
  });
  if (!res.ok) {
    throw new Error(`updateTask failed: ${res.statusText}`);
  }


  const data = await res.json();
  return data.task; // по OpenAPI: handlers.UpdateTaskResponse
}
