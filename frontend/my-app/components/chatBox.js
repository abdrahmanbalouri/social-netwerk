"use client";
import { useState, useRef } from "react";
import "../styles/chat.css";

export default function ChatBox({ user }) {
    const [messages, setMessages] = useState([]);
    const [input, setInput] = useState("");
    const inputRef = useRef(null)
    if (!user) {
        return <div className="loading">Loading user...</div>;
    }

    const emojiArray = [
        "😀", "😃", "😄", "😁", "😆", "😅", "🤣", "😂", "🚀", "💡",
        "😊", "😇", "🙂", "🙃", "😉", "😍", "🥰", "😘", "😗", "😋",
        "😛", "😜", "🤪", "😝", "🤑", "🤗", "🤭", "🤔", "🤨", "😐",
        "😑", "😶", "😏", "😒", "🙄", "😬", "😔", "😪", "🤤", "😴",
        "😷", "🤒", "🤕", "🤢", "🤮", "🥴", "😵", "🤯", "😎", "🤓"
    ];

    const sendMessage = () => {
        if (input.trim() === "") return;
        setMessages([...messages, { text: input, sender: "me" }]);
        setInput("");
    };

    const addEmoji = (emoji) => {
        const cursorPos = inputRef.current.selectionStart;
        const newText =
            input.slice(0, cursorPos) + emoji + input.slice(cursorPos);
        setInput(newText);
        setTimeout(() => inputRef.current.focus(), 0);
    };

    return (
        <div className="chat-container">
            <div className="chat-header">
                <img src={user.image ? user.image : "/uploads/default.png"} alt="Profile" className="profile-pic" />
                <h3 className="username">{user.username}</h3>
            </div>

            <div className="chat-box">
                {messages.length === 0 ? (
                    <p className="no-msg">No messages yet 😶</p>
                ) : (
                    messages.map((msg, index) => (
                        <div key={index} className={`message ${msg.sender}`}>
                            {msg.text}
                        </div>
                    ))
                )}
            </div>

            <div className="emoji-panel">
                {emojiArray.map((emoji, i) => (
                    <span key={i} className="emoji" onClick={() => addEmoji(emoji)}>
                        {emoji}
                    </span>
                ))}
            </div>

            <div className="input-area">
                <input
                    ref={inputRef}
                    type="text"
                    placeholder="Type a message..."
                    value={input}
                    onChange={(e) => setInput(e.target.value)}
                />
                <button onClick={sendMessage}>Send</button>
            </div>
        </div>
    );
}
