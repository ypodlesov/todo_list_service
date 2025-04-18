// src/EditTaskModal.tsx
import React, { useState } from 'react';
import { Task, updateTask } from './api.ts';

interface EditTaskModalProps {
  task: Task;                    // Задача, которую редактируем
  onClose: () => void;           // Закрыть модалку без изменений
  onSave: (updated: Task) => void; // Коллбек, когда мы сохранили задачу
}

/**
 * Простейшая модалка для редактирования title/description
 */
const EditTaskModal: React.FC<EditTaskModalProps> = ({ task, onClose, onSave }) => {
  const [title, setTitle] = useState(task.title);
  const [description, setDescription] = useState(task.description);

  // Сохранить изменения
  const handleSave = async () => {
    try {
      const updated = await updateTask({
        ...task,
        title,
        description,
      });
      // Сообщаем родителю, что задача обновлена
      onSave(updated);
      // Закрываем модалку
      onClose();
    } catch (err) {
      console.error('Failed to update task:', err);
    }
  };

  return (
    <div
      style={{
        position: 'fixed',
        top: 0, left: 0, right: 0, bottom: 0,
        backgroundColor: 'rgba(0,0,0,0.4)',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        zIndex: 9999, // чтобы модалка была поверх всего
      }}
      onClick={onClose} // клик по фону
    >
      <div
        style={{
          background: '#fff',
          padding: 20,
          borderRadius: 8,
          width: 400,
        }}
        onClick={(e) => e.stopPropagation()} // чтобы клик внутри модалки не закрывал её
      >
        <h2>Edit Task</h2>

        <label>Title</label>
        <input
          type="text"
          style={{ display: 'block', width: '100%', marginBottom: 10 }}
          value={title}
          onChange={(e) => setTitle(e.target.value)}
        />

        <label>Description</label>
        <textarea
          style={{ display: 'block', width: '100%', height: 60}}
          value={description}
          onChange={(e) => setDescription(e.target.value)}
        />

        <div style={{ marginTop: 10 }}>
          <button onClick={handleSave}>Save</button>
          <button onClick={onClose} style={{ marginLeft: 10 }}>Cancel</button>
        </div>
      </div>
    </div>
  );
};

export default EditTaskModal;
