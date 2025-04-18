// src/TaskItem.tsx
import React from 'react';
import { DraggableProvided } from 'react-beautiful-dnd';
import { Task } from './api.ts';
import { formatDate } from './utils.ts';
import { MdDragIndicator } from 'react-icons/md';
import './TaskItem.css'; // Импортируйте CSS-файл

interface TaskItemProps {
  task: Task;
  isDraggable: boolean;
  onToggleStatus: (task: Task) => void;
  onClick: () => void;
  provided?: DraggableProvided;
  index?: number;
}

const TaskItem: React.FC<TaskItemProps> = ({
  task,
  isDraggable,
  onToggleStatus,
  onClick,
  provided,
}) => {
  return (
    <div
      ref={provided?.innerRef}
      {...(isDraggable ? provided?.draggableProps : {})}
      className="task-item"
      style={isDraggable ? provided?.draggableProps.style : undefined}
    >
      {/* Чекбокс */}
      <input
        type="checkbox"
        checked={task.status === 2}
        onChange={(e) => {
          e.stopPropagation();
          onToggleStatus(task);
        }}
        className="task-checkbox"
      />

      {/* Основная информация */}
      <div className="task-main" onClick={onClick}>
        <div className="task-title">{task.title}</div>
        <div className="task-date">{formatDate(task.creation_ts)}</div>
      </div>

      {/* Хендл для перетаскивания или заполнение */}
      {isDraggable ? (
        <div
          {...provided?.dragHandleProps}
          className="task-drag-handle"
        >
          <MdDragIndicator size={16} />
        </div>
      ) : (
        <div className="task-drag-placeholder" />
      )}
    </div>
  );
};

export default TaskItem;
