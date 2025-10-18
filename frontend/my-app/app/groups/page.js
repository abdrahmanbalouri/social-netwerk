"use client"
import Navbar from '../../components/Navbar.js';
import { useEffect, useState } from "react";
import "./page.css"
import { useRouter } from 'next/navigation'

import Link from 'next/link'
import styles from './page.css';


export default function () {
    return (
        <>
            <Navbar />
            {/* <AllGroups /> */}
            <GroupsTabs />
        </>
    )
}

function GroupsTabs() {
    const [activeTab, setActiveTab] = useState('my-groups');
    return (
        <>
            <div className='tabs-container '>
                <button className={activeTab === 'my-groups' ? 'tab-button active' : 'tab-button'}
                    onClick={() => setActiveTab('my-groups')}>My Groups</button>
                <button className={activeTab === 'all-groups' ? 'tab-button active' : 'tab-button'}
                    onClick={() => setActiveTab('all-groups')}>All Groups</button>
            </div>

            {/* <div className="group-container"> */}
            {activeTab === 'my-groups' ? (
                <MyGroups />
            ) : (
                <AllGroups />
            )}
            {/* </div> */}
        </>
    );
}

function AllGroups() {
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
        return <div>No group data available.</div>;
    }
    console.log("heeey: ", group)
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


function MyGroups() {
    const [group, setGroup] = useState(null);
    const [loading, setLoading] = useState(true);
    const router = useRouter()
    const handleShow = (group) => {
        console.log("grouuuuuup is :", group);
        localStorage.setItem('selectedGroup', JSON.stringify(group))
        router.push('/groups/posts')
    }



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
        return <div>No group data available.</div>;
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
                        <button onClick={() => handleShow(grp.ID)}>Show</button>
                    </div>
                </div>
            ))}
        </div>
    )
}
