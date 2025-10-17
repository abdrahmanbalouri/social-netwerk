"use client";
import { useState, useRef, useEffect } from "react";
import AddReactionIcon from "@mui/icons-material/AddReaction";
import SendIcon from "@mui/icons-material/Send";
import Link from "next/link";
import "../styles/chat.css";
import { useWS } from "../context/wsContext.js";
import { useChat } from "../context/chatContext.js";

export default function ChatBox({ user }) {
  const [messages, setMessages] = useState([]);
  const [input, setInput] = useState("");
  const [showEmojis, setShowEmojis] = useState(false);
  const inputRef = useRef(null);
  const chatEndRef = useRef(null);
  const { sendMessage, addListener, removeListener } = useWS();
  const { activeChatID, setActiveChatID } = useChat();

  if (!user || user.id == 0) {
    return <div className="loading">Loading user...</div>;
  }

  useEffect(() => {
    setActiveChatID(user.id);
    return () => setActiveChatID(null);
  }, [user.id]);

  setTimeout(() => {
    inputRef.current?.focus();
  }, 0);

  useEffect(() => {
    const fetchMessages = async () => {
      try {
        const response = await fetch(
          `http://localhost:8080/api/getmessages?receiverId=${user.id}`,
          { credentials: "include", method: "GET" }
        );
        if (!response.ok) throw new Error("Failed to fetch messages");
        const data = await response.json();

        if (data.messages) {
          const formattedMessages = data.messages
            .map((msg) => ({              
              text: msg.content,
              sender: msg.senderId === user.id ? "them" : "me",
            }))
            .reverse();

          console.log("Fetched messages:", formattedMessages);
          setMessages(formattedMessages);
        }

      } catch (error) {
        console.error("Error fetching messages:", error);
      }
    };
    fetchMessages();
  }, [user.id]);

  useEffect(() => {
    const handleIncomingMessage = (data) => {
      if (data.from === user.id || data.to === user.id) {
        setMessages((prev) => [
          ...prev,
          { text: data.content, sender: data.from === user.id ? "them" : "me" },
        ]);
      }
    };

    addListener("message", handleIncomingMessage);
    return () => removeListener("message", handleIncomingMessage);
  }, [addListener, removeListener]);

  // 👇 Scroll to bottom whenever messages change
  useEffect(() => {
    if (chatEndRef.current) {
      chatEndRef.current.scrollIntoView({ behavior: "smooth" });
    }
  }, [messages]);

  const emojiArray = [
    "😀",
    "😃",
    "😄",
    "😁",
    "😆",
    "😅",
    "🤣",
    "😂",
    "🚀",
    "💡",
    "😊",
    "😇",
    "🙂",
    "🙃",
    "😉",
    "😍",
    "🥰",
    "😘",
    "😗",
    "😋",
    "😛",
    "😜",
    "🤪",
    "😝",
    "🤑",
    "🤗",
    "🤭",
    "🤔",
    "🤨",
    "😐",
    "😑",
    "😶",
    "😏",
    "😒",
    "🙄",
    "😬",
    "😔",
    "😪",
    "🤤",
    "😴",
    "😷",
    "🤒",
    "🤕",
    "🤢",
    "🤮",
    "🥴",
    "😵",
    "🤯",
    "😎",
    "🤓",
  ];
  const handleSendMessage = () => {
    if (input.trim() === "") return;
    const payload = { receiverId: user.id, messageContent: input, type: "message" };
    sendMessage(payload);
    setInput("");
    setShowEmojis(false);
  };

  const addEmoji = (emoji) => {
    const cursorPos = inputRef.current.selectionStart;
    const newText = input.slice(0, cursorPos) + input.slice(cursorPos) + emoji;
    setInput(newText);
    setTimeout(() => {
      inputRef.current.focus();
      const end = inputRef.current.value.length;
      inputRef.current.setSelectionRange(end, end);
    }, 0);
  };

  return (
    <div className="chat-container">
      <div className="chat-header">
        <Link href={`/profile/${user.id}`}>
          <img
            src={user?.image ? `/uploads/${user.image}` : "/uploads/default.png"}
            alt="user avatar"
          />
        </Link>
        <div className="onlinee">
          <Link href={`/profile/${user.id}`}>
            <span className="username">{user.nickname}</span>
          </Link>
          <span className="on">online</span>
        </div>
      </div>

      {/* Chat box */}
      <div className="chat-box">
        {messages.length === 0 ? (
          <p className="no-msg">No messages yet</p>
        ) : (
          messages.map((msg, index) => (
            <div key={index} className={`message ${msg.sender}`}>
              {msg.text}
            </div>
          ))
        )}
        <div ref={chatEndRef}></div>
      </div>

      {/* Emoji panel */}
      {showEmojis && (
        <div className="emoji-panel">
          {emojiArray.map((emoji, i) => (
            <span key={i} className="emoji" onClick={() => addEmoji(emoji)}>
              {emoji}
            </span>
          ))}
        </div>
      )}

      {/* Input area */}
      <div className="input-area">
        <button className="emoji-toggle" onClick={() => setShowEmojis(!showEmojis)}>
          <AddReactionIcon />
        </button>

        <input
          ref={inputRef}
          type="text"
          placeholder="Type a message..."
          value={input}
          onChange={(e) => setInput(e.target.value)}
          onKeyDown={(e) => e.key === "Enter" && handleSendMessage()}
        />
        <button onClick={handleSendMessage}>
          <SendIcon />
        </button>
      </div>
    </div>
  );
}
