// import Navbar from '../../../components/Navbar.js';
import Navbar from '../../components/Navbar.js';
import React, { useState, useEffect } from 'react';

export function GroupsPage() {
    // 1. Use state to store the fetched group data
    const [group, setGroup] = useState(null);
    const [loading, setLoading] = useState(true);

    // 2. Use the useEffect hook for the side effect (data fetching)
    useEffect(() => {
        fetch('http://localhost:8080/api/groups')
            .then(res => res.json())
            .then(data => {
                setGroup(data); // 3. Store the data in state
                setLoading(false);
            })
            .catch(error => {
                console.error("Failed to fetch group data:", error);
                setLoading(false);
            });
    }, []); // Empty dependency array ensures the fetch runs ONLY once on mount

    // 4. Use JSX to render based on the current state
    if (loading) {
        return <div>Loading groups...</div>;
    }

    if (!group) {
        return <div>No group data available.</div>;
    }

    return (
        <div className="Group-Div">
            {/* Display data from state */}
            {group.title}
        </div>
    );
}


export default function () {
    return (
        <>
            <Navbar />
            <GroupsPage />
        </>
    )
}