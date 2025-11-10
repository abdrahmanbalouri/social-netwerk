"use client";

import { useState, useEffect } from "react";

export default function ShowToast({ message }) {
  const [toast, setToast] = useState(null);

  useEffect(() => {

    if (!message) return;
    setToast(null);
    const timer = setTimeout(() => {
      setToast({ message, id: Date.now() });
    }, 0);
    const autoHide = setTimeout(() => {
      setToast(null);
    }, 3000);

    return () => {
      clearTimeout(timer);
      clearTimeout(autoHide);
    };
  }, [message]);

  return (
    <>
      {toast && (
        <div key={toast.id} className={`toast error`}>
          <span>{toast.message}</span>
          <button onClick={() => setToast(null)} className="toast-close">×</button>
        </div>
      )}
    </>
  );
}
