"use client";
import { useState, useRef } from "react";
import AddReactionIcon from '@mui/icons-material/AddReaction';
import SendIcon from '@mui/icons-material/Send';
import Link from "next/link"
import "../styles/chat.css";

export default function ChatBox({ user }) {
    const [messages, setMessages] = useState([]);
    const [input, setInput] = useState("");
    const [showEmojis, setShowEmojis] = useState(false);
    const inputRef = useRef(null);

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
            </div>

            {showEmojis && (
                <div className="emoji-panel">
                    {emojiArray.map((emoji, i) => (
                        <span
                            key={i}
                            className="emoji"
                            onClick={() => addEmoji(emoji)}
                        >
                            {emoji}
                        </span>
                    ))}
                </div>
            )}

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
                    onKeyDown={(e) => {
                        if (e.key === "Enter") {
                            sendMessage();
                        }
                    }}
                />
                <button onClick={sendMessage}><SendIcon /></button>
            </div>
        </div>
    );
}
