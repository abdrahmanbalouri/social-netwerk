"use client";
import { createContext, useContext, useState } from "react";

const ChatContext = createContext();

export function ChatProvider({ children }) {
  const [activeChatID, setActiveChatID] = useState(null); // null = no chat open

  return (
    <ChatContext.Provider value={{ activeChatID, setActiveChatID }}>
      {children}
    </ChatContext.Provider>
  );
}

export function useChat() {
  return useContext(ChatContext);
}
