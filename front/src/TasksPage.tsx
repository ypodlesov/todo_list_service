// src/TasksPage.tsx
import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  DragDropContext,
  Droppable,
  Draggable,
  DropResult,
} from 'react-beautiful-dnd';

import {
  logout,
  getTasks,
  Task,
  createTask,
  updatePriority,
  updateTask,
} from './api.ts';
import EditTaskModal from './EditTaskModal.tsx';
import CreateTaskModal from './CreateTaskModal.tsx';
import TaskItem from './TaskItem.tsx';
import { DraggableProvided } from 'react-beautiful-dnd';

interface TasksPageProps {
  onLogout: () => void;
}

function reorder(list: Task[], startIndex: number, endIndex: number): Task[] {
  const result = Array.from(list);
  const [removed] = result.splice(startIndex, 1);
  result.splice(endIndex, 0, removed);
  return result;
}

const INT_MAX = 2147483647;
const INT_MIN = -2147483648;

const TasksPage: React.FC<TasksPageProps> = ({ onLogout }) => {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [showCreateModal, setShowCreateModal] = useState(false);

  // For the edit modal
  const [editTask, setEditTask] = useState<Task | null>(null);

  const navigate = useNavigate();

  // 1. Load tasks
  useEffect(() => {
    (async () => {
      try {
        const fetchedTasks = await getTasks();
        setTasks(fetchedTasks);
      } catch (err) {
        console.error(err);
      }
    })();
  }, []);

  // 2. Logout
  const handleLogout = async () => {
    try {
      await logout();
      onLogout();
      navigate('/login');
    } catch (err) {
      console.error(err);
    }
  };

  // 3. Drag & Drop for open tasks
  const handleDragEnd = async (result: DropResult) => {
    const { source, destination } = result;
    if (!destination || destination.index === source.index) {
      return;
    }

    const openTasks = tasks.filter((t) => t.status !== 2);
    const newOrder = reorder(openTasks, source.index, destination.index);

    const closedTasks = tasks.filter((t) => t.status === 2);
    const newAll = [...newOrder, ...closedTasks];
    setTasks(newAll);

    const movedTask = newOrder[destination.index];
    const prevTask = newOrder[destination.index - 1];
    const nextTask = newOrder[destination.index + 1];

    let prevPriority = prevTask?.priority;
    if (prevPriority === undefined) {
      prevPriority = INT_MAX;
    }
    let nextPriority = nextTask?.priority;
    if (nextPriority === undefined) {
      nextPriority = INT_MIN;
    }

    try {
      await updatePriority(movedTask, prevPriority, nextPriority);
    } catch (err) {
      console.error('Failed to update priority:', err);
    }
  };

  // 4. Create a new task (handled via modal)
  const handleCreateTask = async (title: string, description: string) => {
    try {
      const newTask = await createTask(title, description);
      setTasks((prev) => [newTask, ...prev]);
      setShowCreateModal(false);
    } catch (error) {
      console.error('Failed to create task:', error);
    }
  };

  // 5. Toggle task status
  const handleToggleStatus = async (task: Task) => {
    const wasClosed = task.status === 2;
    const newStatus = wasClosed ? 1 : 2;

    try {
      // Update the task status
      const updated = await updateTask({ ...task, status: newStatus });

      // Update the local state
      setTasks((prev) => prev.map((t) => (t.id === updated.id ? updated : t)));

      // If the task is now open, adjust its priority
      if (wasClosed) {
        const openTasks = tasks.filter((t) => t.status !== 2 && t.id !== updated.id);
        const sortedOpen = [...openTasks].sort((a, b) => b.priority - a.priority);

        const topTask = sortedOpen[0];
        let nextPriority = INT_MIN;
        if (topTask) {
          nextPriority = topTask.priority;
        }

        await updatePriority(updated, INT_MAX, nextPriority);

        // Refetch tasks to ensure consistency
        const refetched = await getTasks();
        setTasks(refetched);
      }
    } catch (err) {
      console.error('Failed to update task or priority:', err);
    }
  };

  // 6. Open edit modal
  const handleTaskClick = (task: Task) => {
    if (task.status === 2) return;
    setEditTask(task);
  };

  // 7. Close/save from edit modal
  const handleCloseEdit = () => setEditTask(null);
  const handleSaveEdit = (updated: Task) => {
    setTasks((prev) => prev.map((t) => (t.id === updated.id ? updated : t)));
  };

  const openTasks = tasks.filter((t) => t.status !== 2);
  const closedTasks = tasks.filter((t) => t.status === 2);

  return (
    <div style={{ maxWidth: 600, margin: '20px auto' }}>
      {/* Header and Buttons Container */}
      <div
        style={{
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
          marginBottom: 20,
        }}
      >
        <h1 style={{ margin: 0 }}>Tasks</h1>
        <div>
          <button
            onClick={handleLogout}
            style={{
              backgroundColor: '#007BFF',
              color: '#fff',
              border: 'none',
              padding: '8px 12px',
              borderRadius: 4,
              cursor: 'pointer',
              marginRight: 10,
            }}
          >
            Logout
          </button>
          <button
            onClick={() => setShowCreateModal(true)}
            style={{
              backgroundColor: '#007BFF',
              color: '#fff',
              border: 'none',
              padding: '8px 12px',
              borderRadius: 4,
              cursor: 'pointer',
            }}
          >
            Add Task
          </button>
        </div>
      </div>

      {/* Create Task Modal */}
      {showCreateModal && (
        <CreateTaskModal
          onClose={() => setShowCreateModal(false)}
          onSave={handleCreateTask}
        />
      )}

      {/* Открытые (drag & drop) */}
      <DragDropContext onDragEnd={handleDragEnd}>
        <Droppable droppableId="open-tasks">
          {(provided) => (
            <div ref={provided.innerRef} {...provided.droppableProps}>
              {openTasks.map((task, index) => (
                <Draggable key={task.id} draggableId={String(task.id)} index={index}>
                  {(provided) => (
                    <TaskItem
                      task={task}
                      isDraggable={true}
                      onToggleStatus={handleToggleStatus}
                      onClick={() => handleTaskClick(task)}
                      provided={provided}
                      index={index}
                    />
                  )}
                </Draggable>
              ))}
              {provided.placeholder}
            </div>
          )}
        </Droppable>
      </DragDropContext>

      {openTasks.length > 0 && closedTasks.length > 0 && <hr />}

      {/* Закрытые */}
      <div>
        {closedTasks.map((task) => (
          <TaskItem
            key={task.id}
            task={task}
            isDraggable={false}
            onToggleStatus={handleToggleStatus}
            onClick={() => handleTaskClick(task)}
          />
        ))}
      </div>

      {editTask && (
        <EditTaskModal
          task={editTask}
          onClose={handleCloseEdit}
          onSave={(updated) => {
            handleSaveEdit(updated);
            handleCloseEdit();
          }}
        />
      )}
    </div>
  );
};

export default TasksPage;
