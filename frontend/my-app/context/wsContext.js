"use client";
import { createContext, useContext, useEffect, useRef, useState } from "react";

const WSContext = createContext(null);

export function WSProvider({ children }) {
  const [connected, setConnected] = useState(false);
  const ws = useRef(null);
  const reconnectTimeout = useRef(null);
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
        reconnectTimeout.current = setTimeout(connect, 3000);
      };

      socket.onerror = (error) => {
        console.error("⚠️ WebSocket error:", error);
        socket.close();
      };

      socket.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data);

          // 🔥 Handle online user updates
          // if (data.type === "online_list") {
          //   setOnlineUsers(data.users);
          // } else if (data.type === "status") {
          //   setOnlineUsers((prev) => {
          //     if (data.online) {
          //       return prev.includes(data.userID) ? prev : [...prev, data.userID];
          //     } else {
          //       return prev.filter((id) => id !== data.userID);
          //     }
          //   });
          // }

          // 🔥 Trigger any custom listeners
          if (listeners.current[data.type]) {

            listeners.current[data.type].forEach((cb) => cb(data));
          } else {
            console.warn("⚠️ No listeners for type:", data.type);
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
