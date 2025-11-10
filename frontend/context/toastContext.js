"use clinet"
import { createContext, useContext, useState } from "react";

const ToastContext = createContext(null);

export function ToastProvider({ children }) {
    const [toast, setToast] = useState(null);

    const showToast = (message, type = "error", duration = 3000) => {

        setToast({ message, type });
        setTimeout(() => {
            setToast(null);
        }, duration);
    };


    return <ToastContext.Provider value={{ toast, showToast }}>
        {children}
    </ToastContext.Provider>;
};

export function useToast() {
    return useContext(ToastContext);
}


