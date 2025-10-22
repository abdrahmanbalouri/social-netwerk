"use client"
import Navbar from '../../components/Navbar.js';
import { useEffect, useState } from "react";
import "./page.css"
import { useRouter } from 'next/navigation'
import { GroupCreationTrigger } from '../../components/CreateGroup.js';
import { GroupsTabs } from '../../components/groupTabs.js';


export default function () {
    return (
        <>
            <Navbar />
            {/* <AllGroups /> */}
            <GroupsTabs />
        </>
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
                console.log("data is :", data);
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
    const [group, setGroup] = useState(null);
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
                console.log("data is :", data);
                setGroup(data);
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
                <GroupCreationTrigger />
            </div>
        )
    }
    return (
        <div className="group-container">
            <GroupCreationTrigger />
            {group.map(grp => (
                <div key={grp.ID} className="group-card">
                    <div className="group-content">
                        <h2 className="group-title">{grp.Title}</h2>
                        <p className="group-description">{grp.Description}</p>
                    </div>
                    <div className="group-footer">
                        <button onClick={() => handleShow(grp.ID)}>Show</button>
                    </div>
                </div>
            ))}
        </div>
    )
}

export function createGroup(formData) {
    return fetch('http://localhost:8080/api/groups/add', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify(formData),
    })
        .then(res => {
            if (!res.ok) throw new Error('Failed to create group');
            return res.json();
        })
        .then(data => data)
        .catch(error => {
            console.error('Failed to create new group:', error);
            throw error;
        });
}
