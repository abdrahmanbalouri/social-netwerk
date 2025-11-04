"use client";
import { useState } from "react";
import { AllPosts, Events, GroupChat } from "../app/groups/[id]/page";
import { Calendar, FileText, MessageCircle } from "lucide-react";
import { useParams } from "next/navigation";

export function GroupPostChat() {
  const [activeTab, setActiveTab] = useState("posts");
  const { id } = useParams();
  console.log("---------------", typeof id, id);

  return (
    <div className="group-container">
      <div className="max-width-wrapper">
        {/* Header Tabs */}
        <div className="tabs-container">
          <button
            className={`group-tab-button ${activeTab === "posts" ? "active" : "inactive"
              }`}
            onClick={() => setActiveTab("posts")}
          >
            <FileText className="tab-icon" />
            <span>Posts</span>
          </button>
          <button
            className={`group-tab-button ${activeTab === "chat" ? "active" : "inactive"
              }`}
            onClick={() => setActiveTab("chat")}
          >
            <MessageCircle className="tab-icon" />
            <span>Chat</span>
          </button>
          <button
            className={`group-tab-button ${activeTab === "event" ? "active" : "inactive"
              }`}
            onClick={() => setActiveTab("event")}
          >
            <Calendar className="tab-icon" />
            <span>Events</span>
          </button>
        </div>

        {/* <div className="group-container"> */}
        {
          activeTab === "posts" && (
            <AllPosts grpID={id} />
          )
        }
        {activeTab === "chat" &&

          <GroupChat groupId={id} />
        }
        {
          activeTab === "event" && (
            <Events />
          )
        }
        {/* </div> */}
      </div>
    </div>
  );
}
