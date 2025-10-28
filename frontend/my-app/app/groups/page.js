"use client";
import Navbar from "../../components/Navbar.js";
import { useEffect, useState } from "react";
import "../../styles/groupstyle.css";
import { useRouter } from "next/navigation";
import { GroupCreationTrigger } from "../../components/CreateGroup.js";
import { GroupsTabs } from "../../components/groupTabs.js";
import LeftBar from "../../components/LeftBar.js";
import RightBar from "../../components/RightBar.js";
import { useDarkMode } from "../../context/darkMod.js";
import { Users, ChevronRight } from "lucide-react";
export default function () {
  const { darkMode } = useDarkMode();

  return (
    <div id="div" className={darkMode ? "theme-dark" : "theme-light"}>
      <Navbar />
      {/* main content area */}
      <main className="content" id="contentgroups">
        <LeftBar showSidebar={true} />
        <GroupsTabs />
        <RightBar />
      </main>
    </div>
  );
}

export function AllGroups() {
  const [group, setGroup] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetch("http://localhost:8080/groups", {
      method: "GET",
      credentials: "include",
    })
      .then((res) => res.json())
      .then((data) => {
        console.log("data is :", data);
        setGroup(data);
        setLoading(false);
      })
      .catch((error) => {
        console.error("Failed to fetch group data:", error);
        setLoading(false);
      });
  }, []);
  if (!group) {
    return (
      <div className="empty-state">
        <Users />
        <p className="empty-state-title">No groups found</p>
        <p className="empty-state-text">Check back later for more groups</p>
      </div>
    );
  }
  return (
    <div className="groups-grid">
      {group.map((grp) => (
        <div key={grp.ID} className="group-card">
          <div className="group-header">
            <div className="group-icon-wrapper">
              <div className="group-icon">
                <Users />
              </div>
            </div>
          </div>
          <div className="group-content">
            <h2 className="group-title">{grp.Title}</h2>
            <p className="group-description">{grp.Description}</p>
            <div className="group-footer">
              <div className="group-members">
                <Users />
                <span className="members-text-full">{22} members</span>
                <span className="members-text-short">{220} members</span>
              </div>
              <button className="view-button">
                <span>Join</span>
                <ChevronRight />
              </button>
            </div>
          </div>
        </div>
      ))}
    </div>
  );
}

export function MyGroups() {
  const [group, setGroup] = useState(null);
  const [loading, setLoading] = useState(true);
  const router = useRouter();

  const handleShow = (group) => {
    router.push(`/groups/${group}`);
  };

  useEffect(() => {
    fetch("http://localhost:8080/myGroups", {
      method: "GET",
      credentials: "include",
    })
      .then((res) => res.json())
      .then((data) => {
        console.log("data is :", data);
        setGroup(data);
        setLoading(false);
      })
      .catch((error) => {
        console.error("Failed to fetch group data:", error);
        setLoading(false);
      });
  }, []);
  if (!group) {
    return (
      <div>
        <GroupCreationTrigger />
        <div className="empty-state">
          <Users />
          <p className="empty-state-title">No groups found</p>
          <p className="empty-state-text">
            Create your first group to get started!
          </p>
        </div>
      </div>
    );
  }
  return (
    <>
      <GroupCreationTrigger />
      <div className="groups-grid">
        {group.map((grp) => (
          <div key={grp.ID} className="group-card">
            <div className="group-header">
              <div className="group-icon-wrapper">
                <div className="group-icon">
                  <Users />
                </div>
              </div>
            </div>
            {/* Group Content */}
            <div className="group-content">
              <h2 className="group-title">{grp.Title}</h2>
              <p className="group-description">{grp.Description}</p>
              {/* Group Footer */}
              <div className="group-footer">
                <div className="group-members">
                  <Users />
                  <span className="members-text-full">{22} members</span>
                  <span className="members-text-short">{220}</span>
                </div>
                <button
                  className="view-button"
                  onClick={() => handleShow(grp.ID)}
                >
                  <span>View</span>
                  <ChevronRight />
                </button>
              </div>
            </div>
          </div>
        ))}
      </div>
    </>
  );
}

export async function createGroup(formData) {
  console.log("form data in api call is :", formData);
  return fetch("http://localhost:8080/api/groups/add", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    credentials: "include",
    body: JSON.stringify(formData),
  })
    .then((res) => {
      if (!res.ok) throw new Error("Failed to create group");
      return res.json();
    })
    .then((data) => data)
    .catch((error) => {
      console.error("Failed to create new group:", error);
      throw error;
    });
}
