"use client";
import { createContext, useContext, useEffect, useRef, useState } from "react";

const WSContext = createContext(null);

export function WSProvider({ children }) {
  const [connected, setConnected] = useState(false);
  const ws = useRef(null);
  const reconnectTimeout = useRef(null);

  // 👇 Holds all event listeners, e.g. { message: [cb1, cb2], status: [cb3] }
  const listeners = useRef({});

  // Add a listener for a specific message type
  const addListener = (type, callback) => {
    if (!listeners.current[type]) {
      listeners.current[type] = [];
    }
    listeners.current[type].push(callback);
  };

  // Remove a listener (optional)
  const removeListener = (type, callback) => {
    if (!listeners.current[type]) return;
    listeners.current[type] = listeners.current[type].filter((cb) => cb !== callback);
  };

  // 🔁 Setup & reconnect WebSocket
  useEffect(() => {
    function connect() {
      if (ws.current) {
        console.log("Closing old WebSocket");
        ws.current.close();
      }

      const socket = new WebSocket("ws://localhost:8080/ws");
      ws.current = socket;

      socket.onopen = () => {
        console.log("WebSocket connected");
        setConnected(true);
      };

      socket.onclose = () => {
        console.log("WebSocket closed. Reconnecting in 3s...");
        setConnected(false);
        reconnectTimeout.current = setTimeout(connect, 3000);
      };

      socket.onerror = (error) => {
        console.error("⚠️ WebSocket error:", error);
        socket.close();
      };

      // 🧠 Handle incoming messages and call all relevant listeners
      socket.onmessage = (event) => {
        try {
          console.log("📩 Received:", event);
          const data = JSON.parse(event.data);
          const { type } = data;
          if (listeners.current[type]) {
            listeners.current[type].forEach((cb) => cb(data));
          } else {
            console.warn("No listeners for type:", type, data);
          }
        } catch (err) {
          console.error("❌ Error parsing WebSocket message:", err);
        }
      };
    }

    connect();

    return () => {
      clearTimeout(reconnectTimeout.current);
      ws.current?.close();
    };
  }, []);

  // Send a message easily
  const sendMessage = (msg) => {
    if (ws.current && ws.current.readyState === WebSocket.OPEN) {
      ws.current.send(JSON.stringify(msg));
    } else {
      console.warn("⚠️ WebSocket not connected");
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
