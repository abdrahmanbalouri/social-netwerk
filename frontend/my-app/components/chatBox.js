"use client";
import { useState, useRef, useEffect } from "react";
import AddReactionIcon from "@mui/icons-material/AddReaction";
import SendIcon from "@mui/icons-material/Send";
import Link from "next/link";
import "../styles/chat.css";
import { useWS } from "../context/wsContext.js";
import { useChat } from "../context/chatContext.js";
import { useParams } from "next/navigation";
import { useProfile } from "../context/profile.js";

export default function ChatBox({ user }) {
  const [messages, setMessages] = useState([]);
  const [input, setInput] = useState("");
  const [showEmojis, setShowEmojis] = useState(false);
  const inputRef = useRef(null);
  const chatEndRef = useRef(null);
  const { setActiveChatID } = useChat();
  const [onlineUsers, setonlineUsers] = useState([])
  const { sendMessage, addListener, removeListener } = useWS();
  const id = useParams().id;
  const { Profile } = useProfile();
  const [toast, setToast] = useState(null)
  const showToast = (message, type = "error", duration = 3000) => {
    setToast({ message, type });
    setTimeout(() => {
      setToast(null);
    }, duration);
  };
  if (id == "0" || !id) {
    return (
      <div className="no-chat-selected">
        <div className="no-chat-header">
          <h2>Select a chat to start messaging</h2>
          <br />
          <h4>No chat selected</h4>
        </div>
      </div>
    );
  }
  setTimeout(() => {
    inputRef.current?.focus();
  }, 0);

  useEffect(() => {
    setActiveChatID(user.id);
    return () => setActiveChatID(null);
  }, [user.id]);

  // listen for online/offline status updates
  useEffect(() => {
    const handleOlineUser = (data) => {
      setonlineUsers(data.users)
    }
    sendMessage({ type: "online_list" })
    addListener("online_list", handleOlineUser)
    return () => removeListener("online_list", handleOlineUser)
  }, [addListener, removeListener])

  useEffect(() => {
    const handleLogout = (data) => {
      let useroff = data.userID
      let arr = onlineUsers.filter((id) => {
        return id !== useroff
      })
      setonlineUsers([...arr])
    }

    addListener("logout", handleLogout)
    return () => removeListener("logout", handleLogout)
  }, [addListener, removeListener])




  useEffect(() => {
    const fetchMessages = async () => {
      try {
        const response = await fetch(
          `http://localhost:8080/api/getmessages?receiverId=${user.id}`,
          { credentials: "include", method: "GET" }
        );
        if (!response.ok) throw new Error("Failed to fetch messages");
        const data = await response.json();
console.log("rrrrrrrrrr",data.messages);

        if (data.messages) {
          const formattedMessages = data.messages
            .map((msg) => ({
              text: msg.content,
              sender: msg.senderId === user.id ? "them" : "me",
              time: new Date(msg.createdAt).toLocaleString(),
              name: msg.first_name + " " + msg.last_name,
              image: msg.photo
            }))
            .reverse();

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
          {
            text: data.content,
            sender: data.from === user.id ? "them" : "me",
            time: data.time,
            name: data.name,
            image: data.image
          },
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

  const emojiArray = ["😀", "😃", "😄", "😁", "😆", "😅", "🤣", "😂", "🚀", "💡", "😊", "😇", "🙂", "🙃", "😉", "😍", "🥰", "😘", "😗", "😋", "😛", "😜", "🤪", "😝", "🤑", "🤗", "🤭", "🤔", "🤨", "😐", "😑", "😶", "😏", "😒", "🙄", "😬", "😔", "😪", "🤤", "😴", "😷", "🤒", "🤕", "🤢", "🤮", "🥴", "😵", "🤯", "😎", "🤓"];
  const handleSendMessage = () => {
    if (input.trim() === "") return;
    if (input.length > 1000) {
      showToast("message is too long")
      return
    }
    const payload = {
      receiverId: user.id,
      messageContent: input,
      type: "message",
      name: Profile.first_name,
      image: Profile.image
    };
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
      {toast && (
        <div className={`toast ${toast.type}`}>
          <span>{toast.message}</span>
          <button onClick={() => setToast(null)} className="toast-close">×</button>
        </div>
      )}
        <Link href={`/profile/${user.id}`}>
          <img
            src={
              user?.image ? `/uploads/${user.image}` : "/assets/default.png"
            }
            alt="user avatar"
          />
        </Link>
        <div className="onlinee">
          <Link href={`/profile/${user.id}`}>
            <span className="username">{user.first_name + " " + user.last_name}</span>
          </Link>
          <span className={onlineUsers.includes(user.id) ? "on" : "off"}>
            {onlineUsers.includes(user.id) ? "online" : "offline"}
          </span>
        </div>
      </div>

      {/* Chat box */}
      <div className="chat-box">
        {messages.length === 0 ? (
          <p className="no-msg">No messages yet</p>
        ) : (
          messages.map((msg, index) => (
            <div key={index} className={`message ${msg.sender}`}>
              <div className="message-header">
                <img src={msg?.image ? `/uploads/${msg?.image}` : "/assets/default.png"} alt="profile picture" className="profilePic" />
                <span className="message-name">{msg.name}</span>
              </div>
              <div className="message-bubble">
                <div className="msg-content">{msg.text}</div>
              </div>
              <div className="info-time">
                <span className="time">
                  {new Date(msg.time).toLocaleTimeString("en-US", {
                    hour: 'numeric',
                    minute: '2-digit'
                  })}
                </span>
              </div>
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
        <button
          className="emoji-toggle"
          onClick={() => setShowEmojis(!showEmojis)}
        >
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
