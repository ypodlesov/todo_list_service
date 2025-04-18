import React, { useState } from 'react';
import { Task, updateTask } from './api.ts';

interface EditTaskModalProps {
  task: Task;
  onClose: () => void;
  onSave: (updated: Task) => void;
}

const EditTaskModal: React.FC<EditTaskModalProps> = ({ task, onClose, onSave }) => {
  const [title, setTitle] = useState(task.title);
  const [description, setDescription] = useState(task.description);

  const handleSave = async () => {
    try {
      const updated = await updateTask({
        ...task,
        title,
        description,
      });
      onSave(updated);
      onClose();
    } catch (err) {
      console.error('Failed to update task:', err);
    }
  };

  return (
    <div style={modalOverlayStyle} onClick={onClose}>
      <div style={modalContentStyle} onClick={(e) => e.stopPropagation()}>
        {/* Close Icon */}
        <button onClick={onClose} style={closeButtonStyle} aria-label="Close">
          &times;
        </button>

        {/* Modal Header */}
        <h3 style={modalHeaderStyle}>Edit Task</h3>

        {/* Title Input */}
        <label style={labelStyle}>Title</label>
        <input
          type="text"
          style={inputStyle}
          value={title}
          onChange={(e) => setTitle(e.target.value)}
        />

        {/* Description Input */}
        <label style={labelStyle}>Description</label>
        <textarea
          style={{ ...inputStyle, height: '80px', resize: 'vertical'}}
          value={description}
          onChange={(e) => setDescription(e.target.value)}
        />

        {/* Buttons */}
        <div style={buttonContainerStyle}>
          <button onClick={onClose} style={cancelButtonStyle}>
            Cancel
          </button>
          <button onClick={handleSave} style={saveButtonStyle}>
            Save
          </button>
        </div>
      </div>
    </div>
  );
};

const modalOverlayStyle: React.CSSProperties = {
  position: 'fixed',
  top: 0,
  left: 0,
  right: 0,
  bottom: 0,
  backgroundColor: 'rgba(0, 0, 0, 0.5)',
  display: 'flex',
  justifyContent: 'center',
  alignItems: 'center',
  zIndex: 1000,
};

const modalContentStyle: React.CSSProperties = {
  backgroundColor: '#fff',
  padding: '24px',
  borderRadius: '12px',
  width: '90%',
  maxWidth: '600px',
  boxShadow: '0 4px 10px rgba(0, 0, 0, 0.3)',
  position: 'relative',
};

const closeButtonStyle: React.CSSProperties = {
  position: 'absolute',
  top: '12px',
  right: '12px',
  background: 'none',
  border: 'none',
  fontSize: '20px',
  cursor: 'pointer',
  color: '#666',
};

const modalHeaderStyle: React.CSSProperties = {
  margin: '0 0 20px',
  fontSize: '20px',
  fontWeight: 'bold',
  textAlign: 'left',
};

const labelStyle: React.CSSProperties = {
  display: 'block',
  marginBottom: '8px',
  fontSize: '14px',
  fontWeight: 'bold',
};

const inputStyle: React.CSSProperties = {
  width: '100%',
  padding: '8px 12px',
  border: '1px solid #ccc',
  borderRadius: '6px',
  fontSize: '14px',
  marginBottom: '16px',
  boxSizing: 'border-box',
};

const buttonContainerStyle: React.CSSProperties = {
  display: 'flex',
  justifyContent: 'flex-end',
  gap: '10px',
};

const cancelButtonStyle: React.CSSProperties = {
  backgroundColor: '#fff',
  color: '#000',
  border: '1px solid #ccc',
  borderRadius: '6px',
  padding: '8px 16px',
  cursor: 'pointer',
  fontSize: '14px',
};

const saveButtonStyle: React.CSSProperties = {
  backgroundColor: '#007BFF',
  color: '#fff',
  border: 'none',
  borderRadius: '6px',
  padding: '8px 16px',
  cursor: 'pointer',
  fontSize: '14px',
};

export default EditTaskModal;
