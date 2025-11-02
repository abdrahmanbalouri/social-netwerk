"use client";
import { useState } from "react";
import "../styles/groupstyle.css";
import { MyGroups, AllGroups } from "../app/groups/page";
import { Globe, Users } from "lucide-react";
export function GroupsTabs() {
    const [activeTab, setActiveTab] = useState("myGroups");
    return (
        <div className="group-container">
            <div className="max-width-wrapper">
                {/* Header Tabs */}
                <div className="tabs-container">
                    <button
                        onClick={() => setActiveTab("myGroups")}
                        className={`group-tab-button ${activeTab === "myGroups" ? "active" : "inactive"
                            }`}
                    >
                        <Users className="tab-icon" />
                        <span className="tab-text-full">My Groups</span>
                        <span className="tab-text-short">My</span>
                    </button>
                    <button
                        onClick={() => setActiveTab("allGroups")}
                        className={`group-tab-button ${activeTab === "allGroups" ? "active" : "inactive"
                            }`}
                    >
                        <Globe className="tab-icon" />
                        <span className="tab-text-full">All Groups</span>
                        <span className="tab-text-short">All</span>
                    </button>
                </div>
                {activeTab === "myGroups" ? <MyGroups /> : <AllGroups />}
            </div>
        </div>
    );
}