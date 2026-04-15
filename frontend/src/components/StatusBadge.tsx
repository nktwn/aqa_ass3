// src/components/StatusBadge.tsx
import React from 'react';

export function getStatusBadge(status: string): React.ReactElement {
    const baseClass = "px-2 py-1 rounded text-xs font-semibold";

    switch (status) {
        case 'Pending':
            return <span className={`bg-yellow-100 text-yellow-800 ${baseClass}`}>🕒 Ожидает</span>;
        case 'In Progress':
            return <span className={`bg-blue-100 text-blue-800 ${baseClass}`}>🔄 В процессе</span>;
        case 'Completed':
            return <span className={`bg-green-100 text-green-800 ${baseClass}`}>✅ Завершён</span>;
        case 'Cancelled':
            return <span className={`bg-red-100 text-red-800 ${baseClass}`}>❌ Отменён</span>;
        default:
            return <span className={`bg-gray-100 text-gray-800 ${baseClass}`}>{status}</span>;
    }
}
