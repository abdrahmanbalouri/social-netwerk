"use client";
import Navbar from "../../../components/Navbar.js";
import { GroupPostChat } from "../../../components/groupPostCat.js";
import Post from "../../../components/Post.js";
import { useEffect, useRef, useState } from "react";
import "../../../styles/groupstyle.css";
import { useParams } from "next/navigation";
import { PostCreationTrigger } from "../../../components/cretaePostGroup.js";
import LeftBar from "../../../components/LeftBar.js";
import RightBarGroup from "../../../components/RightBarGroups.js";
import { useDarkMode } from "../../../context/darkMod.js";
import { FileText, SendIcon } from "lucide-react";
import AddReactionIcon from "@mui/icons-material/AddReaction";
import { useWS } from "../../../context/wsContext.js";
import "../../../styles/chat.css";

export default function () {
  const { darkMode } = useDarkMode();

  return (
    <div className={darkMode ? "theme-dark" : "theme-light"}>
      <Navbar />
      <main className="content">
        <LeftBar showSidebar={true} />
        {/* <AllPosts /> */}
        <GroupPostChat />
        {/* <RightBar /> */}
        <RightBarGroup onClick={sendRequest} />
      </main>
      {/* <AllPosts /> */}
      {/* <CreatePost /> */}
    </div>
  );
}

function sendRequest(invitedUserID, grpID) {
  console.log("invited user id is : ", invitedUserID);
  console.log("group id is : ", grpID);

  fetch(`http://localhost:8080/group/invitation/${grpID}`, {
    method: "POST",
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      invitedUsers: [invitedUserID],
    }),
  })
    .then(async (res) => {
      const temp = await res.json();
      console.log("temp is :", temp);
    })
    .catch((error) => {
      console.log("error sending invitation to the user :", error);
    });
}

export function AllPosts({ grpID }) {
  const [posts, setPost] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!grpID) {
      setLoading(false);
      return;
    }

    fetch(`http://localhost:8080/group/fetchPosts/${grpID}`, {
      method: "GET",
      credentials: "include",
    })
      .then((res) => {
        if (!res.ok) throw new Error("Failed to fetch posts");
        return res.json();
      })
      .then((data) => {
        setPost(data || []);
        setLoading(false);
      })
      .catch((error) => {
        console.error("Failed to fetch posts:", error);
        setLoading(false);
      });
  }, []);
  // if (!posts) {
  //     return (
  //         <div>
  //             <PostCreationTrigger setPost = {setPost}/>
  //             <div>There is no post yeeeeeet.</div>
  //         </div>
  //     );
  // }

  console.log("posts are :", posts);
  return (
    <div>
      <PostCreationTrigger setPost={setPost} />
      <div className="content-area">
        {posts.length === 0 ? (
          <div className="group-empty-state">
            <FileText />
            <p className="group-empty-state-text">
              No posts yet. Be the first to share something!
            </p>
          </div>
        ) : (
          posts.map((post) => (
            <Post
              key={post.id}
              post={post}
              onGetComments={GetComments}
              ondolike={AddLike}
            />
          ))
        )}
      </div>
    </div>
  );
}
export function LastPost() {
  const [posts, setPosts] = useState([]);
  const [loading, setLoading] = useState(true);
  const params = useParams();
  const grpID = params.id;

  useEffect(() => {
    if (!grpID) {
      setLoading(false);
      return;
    }

    fetch(`http://localhost:8080/group/fetchPost/${grpID}`, {
      method: "GET",
      credentials: "include",
    })
      .then((res) => {
        if (!res.ok) throw new Error("Failed to fetch last post");
        return res.json();
      })
      .then((data) => {
        setPosts(data);
        setLoading(false);
      })
      .catch((error) => {
        console.error("Failed to fetch last post:", error);
        setLoading(false);
      });
  }, [grpID]);

  if (loading) {
    return <div>Loading...</div>;
  }

  if (!posts || posts.length === 0) {
    return (
      <div>
        <PostCreationTrigger setPosts={setPosts} />
        <div>There is no post yet.</div>
      </div>
    );
  }

  console.log("posts are:", posts);

  return (
    <div>
      <PostCreationTrigger setPosts={setPosts} />
      <div className="content-area">
        {posts.map((post) => (
          <Post
            key={post.id}
            post={post}
            onGetComments={GetComments}
            ondolike={AddLike}
          />
        ))}
      </div>
    </div>
  );
}

function GetComments(post) {
  // setSelectedPost({
  //     id: post.id,
  //     title: post.title || post.post_title || "Post"
  // });
}

function AddLike() {}

export async function CreatePost(groupId, formData) {
  const data = new FormData();
  data.append("postData", JSON.stringify(formData));

  const res = await fetch(`http://localhost:8080/group/addPost/${groupId}`, {
    method: "POST",
    credentials: "include",
    body: data,
  });

  const text = await res.text();

  if (!res.ok) throw new Error("Failed to create a post for groups");

  return JSON.parse(text);
}

// ===========================
// group chat component
// ===========================

export function GroupChat({ groupId }) {
  console.log("group chat component rendered");
  const [id, setId] = useState(null);
  const [messages, setMessages] = useState([]);
  const [input, setInput] = useState("");
  const [showEmojis, setShowEmojis] = useState(false);
  const inputRef = useRef(null);
  const chatEndRef = useRef(null);
  const { sendMessage, addListener, removeListener } = useWS();
  useEffect(() => {
    fetch("http://localhost:8080/api/me", {
      credentials: "include",
      method: "GET",
    })
      .then((res) => res.json())
      .then((data) => {
      console.log("user id in group chat------------ : ", data.user_id);
        setId(data.user_id)
      })
      .catch((err) => console.error(err));
  });

  setTimeout(() => {
    inputRef.current?.focus();
  }, 0);

  useEffect(() => {
    const fetchMessages = async () => {
      try {
        const response = await fetch(
          `http://localhost:8080/api/getGroupMessages?groupId=${groupId}`,
          { credentials: "include", method: "GET" }
        );
        if (!response.ok) throw new Error("Failed to fetch messages");
        const data = await response.json();
        if (data.messages) {
          console.log("user id in group chat message listener23131 : ", id);

          const formattedMessages = data.messages
            .map((msg) => ({
              text: msg.content,
              sender: msg.senderId !== id ? "them" : "me",
              time: new Date(msg.createdAt).toLocaleString(),
            }))
            .reverse();

          setMessages(formattedMessages);
        }
      } catch (error) {
        console.error("Error fetching messages:", error);
      }
    };
    fetchMessages();
  }, [id]);

  useEffect(() => {
    const handleIncomingMessage = (data) => {
      console.log("user id in group chat message listener : ", id);
      setMessages((prev) => [
        ...prev,
        {
          text: data.content,
          sender: data.senderId !== id ? "them" : "me",
          time: data.time,
        },
      ]);
    };

    addListener("group_message", handleIncomingMessage);
    return () => removeListener("group_message", handleIncomingMessage);
  }, [addListener, removeListener, id]);

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

  function handleSendGroupChatMessage() {
    console.log("=========", typeof groupId);
    if (input.trim() === "") return;
    const payload = {
      groupID: groupId,
      messageContent: input,
      type: "group_message",
    };
    sendMessage(payload);
    setInput("");
    setShowEmojis(false);
  }

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
    <div className="group-chat-container">
      {/* Chat box */}
      <div className="chat-box">
        {messages.length === 0 ? (
          <p className="no-msg">No messages yet</p>
        ) : (
          messages.map((msg, index) => (
            <div key={index} className={`message ${msg.sender}`}>
              <div className="msg-content">{msg.text}</div>
              <div className="info-time">
                <span className="time">
                  {new Date(msg.time).toLocaleTimeString("en-US")}
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
          onKeyDown={(e) => e.key === "Enter" && handleSendGroupChatMessage()}
        />
        <button onClick={handleSendGroupChatMessage}>
          <SendIcon />
        </button>
      </div>
    </div>
  );
}
