"use client"
import Navbar from '../../components/Navbar.js';
import { GroupCard } from '../../components/groupCard.js';
import { useEffect, useState } from "react";
import "../../styles/groupstyle.css"
import { useRouter, useSelectedLayoutSegments } from 'next/navigation'
import { GroupCreationTrigger } from '../../components/CreateGroup.js';
import { GroupsTabs } from '../../components/groupTabs.js';
import LeftBar from '../../components/LeftBar.js';
import RightBar from '../../components/RightBar.js';
import { useDarkMode } from '../../context/darkMod.js';

export default function () {
    const { darkMode } = useDarkMode();

    return (
        <div id="div" className={darkMode ? 'theme-dark' : 'theme-light'}>
            <Navbar />
            {/* main content area */}
            <main className="content" id="contentgroups">
                <LeftBar showSidebar={true} />
                <GroupsTabs />
                <RightBar />
            </main>
        </div>
    )
}

export function AllGroups() {
    const [group, setGroup] = useState(null);
    const [loading, setLoading] = useState(true);


    useEffect(() => {
        fetch('http://localhost:8080/groups', {
            method: 'GET',
            credentials: 'include',
        })
            .then(res => res.json())
            .then(data => {
                console.log("All the groups are :", data);
                setGroup(data);
                setLoading(false);
            })
            .catch(error => {
                console.error("Failed to fetch group data:", error);
                setLoading(false);
            });
    }, [])
    if (!group) {
        return <div>You don’t have any groups yet. Go create One ! </div>;
    }
    return (
        <div className="group-container">
            {group.map(grp => (
                <div key={grp.ID} className="group-card">
                    <div className="group-content">
                        <h2 className="group-title">{grp.Title}</h2>
                        <p className="group-description">{grp.Description}</p>
                    </div>
                    <div className="group-footer">
                        <button className="join-button">Join</button>
                    </div>
                </div>
            ))}
        </div>
    )
}

export function MyGroups() {
    console.log("MyGroups rendered");
    const [group, setGroup] = useState([]);
    const [loading, setLoading] = useState(true);
    const router = useRouter()

    const handleShow = (group) => {
        router.push(`/groups/${group}`);
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
    if (!group) {
        return (
            <div>
                <GroupCreationTrigger setGroup={setGroup} />
            </div>
        )
    }
    console.log("groups after are :", group);
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

export function createGroup(formData) {
    console.log("inside Create Group function");
    return fetch('http://localhost:8080/api/groups/add', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify(formData),
    })
        .then(async res => {
            console.log("form data ha s: ", formData);
            if (!res.ok) throw new Error('Failed to create group');
            // console.log("result :", res);
            // console.log("result  :",await res.text());
            const groupIS = await res.json()
            // console.log("new group is :", groupIS);
            const SendInvitations = await fetch(
                `http://localhost:8080/group/invitation/${groupIS.ID}`,
                {
                    method: 'POST',
                    credentials: 'include',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({
                        invitedUsers: formData.invitedUsers,
                    }),
                }
            );

            console.log("SendInvitations ::", await SendInvitations.json());
            return groupIS
        })
        // .then(createdGroup => { return createdGroup })
        .catch(error => {
            console.error('Failed to create new group:', error);
            throw error;
        });
}
