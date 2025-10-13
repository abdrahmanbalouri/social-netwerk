"use client";
import { createContext, useContext, useEffect, useRef, useState } from "react";

const WSContext = createContext(null);

export function WSProvider({ children }) {
    const [connected, setConnected] = useState(false);
    const ws = useRef(null);

    useEffect(() => {
        ws.current = new WebSocket("ws://localhost:8080/ws");

        ws.current.onopen = () => setConnected(true);
        ws.current.onclose = () => setConnected(false);

        return () => ws.current.close();
    }, []);

    return (
        <WSContext.Provider value={{ ws: ws.current, connected }}>
            {children}
        </WSContext.Provider>
    );
}

export function useWS() {
    return useContext(WSContext);
}
