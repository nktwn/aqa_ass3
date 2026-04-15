'use client';

import { useState } from 'react';

export default function SignatureModal({
                                           isOpen,
                                           onClose,
                                           onConfirm,
                                       }: {
    isOpen: boolean;
    onClose: () => void;
    onConfirm: (signatureCode: string) => void;
}) {
    const [signatureCode, setSignatureCode] = useState('');
    const [error, setError] = useState('');

    const handleSubmit = () => {
        if (!signatureCode.trim()) {
            setError('Введите код подписи');
            return;
        }
        setError('');
        onConfirm(signatureCode);
    };

    if (!isOpen) return null;

    return (
        <div className="fixed inset-0 bg-black/50 z-50 flex justify-center items-center">
            <div className="bg-white p-6 rounded-xl shadow-lg w-[90%] max-w-sm space-y-4">
                <h2 className="text-xl font-semibold text-center">🔐 Подтверждение подписи</h2>
                <p className="text-sm text-gray-600 text-center">Введите ваш код</p>

                <input
                    type="password"
                    placeholder="Ваш код подписи"
                    value={signatureCode}
                    onChange={(e) => setSignatureCode(e.target.value)}
                    className="w-full border rounded px-4 py-2"
                />

                {error && <p className="text-red-600 text-sm">{error}</p>}

                <div className="flex justify-end gap-2">
                    <button onClick={onClose} className="btn-outline-primary">Отмена</button>
                    <button onClick={handleSubmit} className="btn-primary">Подписать</button>
                </div>
            </div>
        </div>
    );
}
