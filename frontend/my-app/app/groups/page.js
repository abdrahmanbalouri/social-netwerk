"use client"
import Navbar from '../../components/Navbar.js';
import { useEffect, useState } from "react";
import styles from './page.module.css';


export default function () {
    return (
        <>
            <Navbar />
            <Groups />
        </>
    )
}

function Groups() {
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
        <div className={styles["group-container"]}>
            {group.map(grp => (
                <div key={grp.ID} className={styles["group-card"]}>
                    <div className={styles["group-content"]}>
                        <h2 className={styles["group-title"]}>{grp.Title}</h2>
                        <p className={styles["group-description"]}>{grp.Description}</p>
                    </div>
                    <div className={styles["group-footer"]}>
                        <button className={styles["join-button"]}>Join</button>
                    </div>
                </div>
            ))}
        </div>
    )
}

