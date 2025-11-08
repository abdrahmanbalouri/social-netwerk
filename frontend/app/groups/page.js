"use client";
import Navbar from "../../components/Navbar.js";
import { useEffect, useState } from "react";
import "../../styles/groupstyle.css";
import { useRouter } from "next/navigation";
import { GroupCard } from "../../components/groupCard.js";
import { GroupCreationTrigger } from "../../components/CreateGroup.js";
import { GroupsTabs } from "../../components/groupTabs.js";
import LeftBar from "../../components/LeftBar.js";
import RightBar from "../../components/RightBar.js";
import { useDarkMode } from "../../context/darkMod.js";
import { Users, ChevronRight } from "lucide-react";
import { Toaster, toast } from "sonner"
import { useWS } from "../../context/wsContext.js";
import { useProfile } from "../../context/profile.js";
import { send } from "process";

// import RightBarGroup from '../../components/RightBarGroups.js';

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

async function JoinGroup(grpID, setJoining) {

  try {
    const res = await fetch(`http://localhost:8080/group/invitation/${grpID}`, {
      method: 'POST',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        InvitationType: "join",
        invitedUsers: [],
      }),
    });

    const temp = await res.json();
    
    return temp;
  } catch (error) {
    console.error("error sending invitation to join the group :", error);
    return { error: error.message };
  } finally {
    setJoining(false);
  }
}

export function AllGroups() {
  const [group, setGroup] = useState(null);
  const [loading, setLoading] = useState(true);
  const [joining, setJoining] = useState(false)
  const { sendMessage } = useWS();
  const { Profile } = useProfile();


  useEffect(() => {
    fetch("http://localhost:8080/groups", {
      method: "GET",
      credentials: "include",
    })
      .then((res) => res.json())
      .then((data) => {
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
      <div className="group-empty-state">
        <Users />
        <p className="group-empty-state-title">No groups found</p>
        <p className="group-empty-state-text">
          Check back later for more groups
        </p>
      </div>
    );
  }
  return (
    <div className="groups-grid">
      <Toaster position="bottom-right" richColors />
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
                <span className="members-text-full">{grp.MemberCount} members</span>
                <span className="members-text-short">{grp.MemberCount} members</span>
              </div>
              <button className="view-button" onClick={() => {
                toast.success("Join request sent!");
                JoinGroup(grp.ID, setJoining)
                sendMessage({
                  type: "joinRequest",
                  from: Profile.id,
                  receiverId: grp.ID,
                  messageContent: "",
                });
              }}>
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
  const [group, setGroup] = useState([]);
  const [loading, setLoading] = useState(true);
  const router = useRouter()

  const handleShow = (groupId) => {
    router.push(`/groups/${groupId}`);
  };

  useEffect(() => {
    fetch('http://localhost:8080/myGroups', {
      method: 'GET',
      credentials: 'include',
    })
      .then(res => res.json())
      .then(data => {
        setGroup(data || []);
        setLoading(false);
      })
      .catch(error => {
        console.error("Failed to fetch group data:", error);
        setLoading(false);
      });
  }, [])

  return (
    <div className="group-container">
      <GroupCreationTrigger
        setGroup={setGroup}
      />
      {group.map(grp => (
        <GroupCard
          key={grp.ID}
          group={grp}
          onShow={handleShow}
        />
      ))}
    </div>
  )
}

export async function createGroup(formData) {
  console.log("sending request...");
  console.log("------------------",JSON.stringify(formData));
  return (
    fetch("http://localhost:8080/api/groups/add", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      credentials: "include",
      body: JSON.stringify(formData),
    })
      .then(async (res) => {
        if (!res.ok) throw new Error("Failed to create group :");
        const groupIS = await res.json();
        return groupIS;
      })
      // .then(createdGroup => { return createdGroup })
      .catch((error) => {
        console.error("Failed to create new group:", error);
        throw error;
      })
  );
}