"use client";
import { createContext, useContext, useEffect, useState } from "react";

const WSContext = createContext(null);

export function WSProvider({ children }) {
  const [ws, setWs] = useState(null);
  const [connected, setConnected] = useState(false);

  useEffect(() => {
    const socket = new WebSocket("ws://localhost:8080/ws");
    setWs(socket);

    socket.onopen = () => setConnected(true);
    socket.onclose = () => setConnected(false);

    return () => socket.close();
  }, []);

  return (
    <WSContext.Provider value={{ ws, connected }}>
      {children}
    </WSContext.Provider>
  );
}

export function useWS() {
  return useContext(WSContext);
}
