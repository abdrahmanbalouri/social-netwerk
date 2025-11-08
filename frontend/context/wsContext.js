"use client";
import { createContext, useContext, useEffect, useRef, useState } from "react";
import { useRouter } from "next/navigation";
const WSContext = createContext(null);

export function WSProvider({ children }) {
  const route = useRouter();
  const [connected, setConnected] = useState(false);
  const ws = useRef(null);
  const listeners = useRef({});

  const addListener = (type, callback) => {
    if (!listeners.current[type]) listeners.current[type] = [];
    listeners.current[type].push(callback);
  };

  const removeListener = (type, callback) => {
    if (!listeners.current[type]) return;
    listeners.current[type] = listeners.current[type].filter((cb) => cb !== callback);
  };

  useEffect(() => {
    function connect() {
      if (ws.current) {
        ws.current.close();
      }

      const socket = new WebSocket("ws://localhost:8080/ws");
      ws.current = socket;

      socket.onopen = () => {
        setConnected(true);
      };

      socket.onclose = () => {
        setConnected(false);
      };

      socket.onerror = () => {
        socket.close();
      };

      socket.onmessage = async (event) => {
        try {
          const res = await fetch("http://localhost:8080/api/me", {
            method: "GET",
            credentials: "include",
          })
          if (!res.ok) {
            ws.current.close();
            return;
          }

          const data = JSON.parse(event.data);
          // 🔥 Trigger any custom listeners
          if (listeners.current[data.type]) {

            listeners.current[data.type].forEach((cb) => cb(data));
          }
        } catch (err) {
          console.error("❌ Error parsing WebSocket message:", err);
        }
      };
    }

    connect();

    return () => {
      ws.current?.close();
    };
  }, []);

  const sendMessage = async (msg) => {
    try {
      const res = await fetch("http://localhost:8080/api/me", {
        method: "GET",
        credentials: "include",
      })
      if (!res.ok) {
        ws.current.close();
        return;
      }
      
      if (ws.current && ws.current.readyState === WebSocket.OPEN) {
        ws.current.send(JSON.stringify(msg));
      }

    } catch (err) {
      console.error("❌ Error sending WebSocket message:", err);
    }
  };

  return (
    <WSContext.Provider
      value={{
        ws: ws.current,
        connected,
        sendMessage,
        addListener,
        removeListener,
      }}
    >
      {children}
    </WSContext.Provider>
  );
}

export function useWS() {
  return useContext(WSContext);
}
