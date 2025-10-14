"use client";
import { useState, useRef, useEffect } from "react";
import AddReactionIcon from '@mui/icons-material/AddReaction';
import SendIcon from '@mui/icons-material/Send';
import Link from "next/link"
import "../styles/chat.css";
import { useWS } from "../context/wsContext.js";


export default function ChatBox({ user }) {
    const [messages, setMessages] = useState([]);
    const [input, setInput] = useState("");
    const [showEmojis, setShowEmojis] = useState(false);
    const inputRef = useRef(null);
    const { ws, connected } = useWS();
    console.log("user in chatbox:", ws, "connection", connected);

    if (!user) {
        return <div className="loading">Loading user...</div>;
    }
    setTimeout(() => {
        inputRef.current.focus();
    }, 0)

    useEffect(() => {
        if (!ws) return;

        ws.onmessage = (event) => {
            console.log("📩 Received:", event.data);
            const msg = JSON.parse(event.data);

            if (msg.type === "message") {
                setMessages((prev) => [...prev, {
                    text: msg.content,
                    sender: msg.from === user.id ? "them" : "me",
                }]);
            }
        };

        return () => {
            ws.onmessage = null;
        };
    }, [ws]);



    const emojiArray = [
        "😀", "😃", "😄", "😁", "😆", "😅", "🤣", "😂", "🚀", "💡",
        "😊", "😇", "🙂", "🙃", "😉", "😍", "🥰", "😘", "😗", "😋",
        "😛", "😜", "🤪", "😝", "🤑", "🤗", "🤭", "🤔", "🤨", "😐",
        "😑", "😶", "😏", "😒", "🙄", "😬", "😔", "😪", "🤤", "😴",
        "😷", "🤒", "🤕", "🤢", "🤮", "🥴", "😵", "🤯", "😎", "🤓"
    ];

    const sendMessage = () => {
        if (input.trim() === "") return;
        const payload = {
            receiverId: user.id,
            messageContent: input,
            type: "message",
        };

        if (connected && ws) {
            ws.send(JSON.stringify(payload));
        }
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
