import React from 'react';
import './GenreBadge.css';

interface GenreBadgeProps {
  name: string;
  onClick?: () => void;
}

export const GenreBadge: React.FC<GenreBadgeProps> = ({ name, onClick }) => {
  return (
    <span className={`genre-badge ${onClick ? 'clickable' : ''}`} onClick={onClick}>
      {name}
    </span>
  );
};

