"use client";

import { useState, useEffect } from "react";

export default function ShowToast({ message }) {
  const [toast, setToast] = useState(null);

  useEffect(() => {
   // if (!message) return;

    setToast({ message });

    const timer = setTimeout(() => {
      setToast(null);
    }, 3000);

    return () => clearTimeout(timer);
  }, [message]);

  return (
    <>
      {toast && (
        <div className={`toast error`}>
          <span>{toast.message}</span>
          <button onClick={() => setToast(null)} className="toast-close">×</button>
        </div>
      )}
    </>
  );
}
