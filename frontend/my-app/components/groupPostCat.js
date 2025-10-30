import { useState } from "react";
import { AllPosts } from "../app/groups/[id]/page";
import { FileText, MessageCircle } from "lucide-react";

export function GroupPostChat() {
  const [activeTab, setActiveTab] = useState("posts");
  return (
    <div className="group-container">
      <div className="max-width-wrapper">
        {/* Header Tabs */}
        <div className="tabs-container">
          <button
            className={`group-tab-button ${
              activeTab === "posts" ? "active" : "inactive"
            }`}
            onClick={() => setActiveTab("posts")}
          >
            <FileText className="tab-icon" />
            <span>Posts</span>
          </button>
          <button
            className={`group-tab-button ${
              activeTab === "chat" ? "active" : "inactive"
            }`}
            onClick={() => setActiveTab("chat")}
          >
            <MessageCircle className="tab-icon" />
            <span>Chat</span>
          </button>
        </div>

        {/* <div className="group-container"> */}
        {activeTab === "posts" ? (
          <AllPosts />
        ) : (
          <div className="group-empty-state">
            <MessageCircle className="tab-icon" />
            <p className="group-empty-state-text">Chat interface coming soon</p>
          </div>
        )}
        {/* </div> */}
      </div>
    </div>
  );
}
