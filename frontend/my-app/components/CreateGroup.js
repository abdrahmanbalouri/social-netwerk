import { useState, useEffect } from 'react';
import "../styles/createGroup.css"
import { createGroup } from '../app/groups/page';

export function CreateGroupForm({ users, onSubmit, onCancel }) {
    console.log("useeeeers are :", users);
    const [groupTitle, setGroupTitle] = useState('');
    const [groupDescription, setGroupDescription] = useState('');
    const [searchQuery, setSearchQuery] = useState('');
    const [selectedUsers, setSelectedUsers] = useState([]);
    const [showSuggestions, setShowSuggestions] = useState(false);

    // Filter users based on search query
    // useEffect(() => {
    //     if (searchQuery.trim()) {
    //         const filtered = users.filter(user =>
    //             user.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
    //             user.username.toLowerCase().includes(searchQuery.toLowerCase())
    //         );
    //         setFilteredUsers(filtered);
    //     } else {
    //         setFilteredUsers([]);
    //     }
    // }, [searchQuery, users]);

    const handleUserSelect = (user) => {
        if (!selectedUsers.find(u => u.id === user.id)) {
            setSelectedUsers([...selectedUsers, user]);
            setSearchQuery('');
        }
    };

    const handleUserRemove = (userId) => {
        setSelectedUsers(selectedUsers.filter(u => u.id !== userId));
    };

    const handleSubmit = (e) => {
        e.preventDefault();

        const groupData = {
            title: groupTitle,
            description: groupDescription,
            invitedUsers: selectedUsers.map(u => u.id),
        };

        onSubmit(groupData);
        onSubmit(selectedUsers);

        // Reset form after submission
        setGroupTitle('');
        setGroupDescription('');
        setSelectedUsers([]);
        setSearchQuery('');
    };

    console.log("filterd users areeeee :", users);

    return (
        <div className="create-group-form-container">
            <div className="create-group-form-card">
                <h1 className="form-title">Create a New Group</h1>

                <form onSubmit={handleSubmit}>
                    {/* Group Title */}
                    <div className="form-field">
                        <label htmlFor="groupTitle" className="form-label">Group Title</label>
                        <input
                            id="groupTitle"
                            type="text"
                            value={groupTitle}
                            onChange={(e) => setGroupTitle(e.target.value)}
                            placeholder="Enter group title"
                            className="form-input"
                            required
                        />
                    </div>

                    {/* Group Description */}
                    <div className="form-field">
                        <label htmlFor="groupDescription" className="form-label">Description</label>
                        <textarea
                            id="groupDescription"
                            value={groupDescription}
                            onChange={(e) => setGroupDescription(e.target.value)}
                            placeholder="What's your group about?"
                            className="form-textarea"
                            rows="4"
                            required
                        />
                    </div>

                    {/* Invite Users */}
                    <div className="form-field">
                        <label htmlFor="inviteUsers" className="form-label">Invite Members</label>

                        {/* Selected Users */}
                        {selectedUsers.length > 0 && (
                            <div className="selected-users">
                                {selectedUsers.map(user => (
                                    <div key={user.id} className="selected-user-tag">
                                        <span className="user-avatar-small">
                                            {user.nickname.charAt(0).toUpperCase()}
                                        </span>
                                        <span>{user.name}</span>
                                        <button
                                            type="button"
                                            onClick={() => handleUserRemove(user.id)}
                                            className="remove-user-btn"
                                        >
                                            ✕
                                        </button>
                                    </div>
                                ))}
                            </div>
                        )}

                        {/* Search Input */}
                        <div className="search-container">
                            <input
                                id="inviteUsers"
                                type="text"
                                value={searchQuery}
                                onClick={() => setShowSuggestions(prev => !prev)} 
                                onChange={(e) => setSearchQuery(e.target.value)}
                                placeholder="Search users to invite..."
                                className="form-input search-input"
                            />

                            {/* User Suggestions */}
                            {showSuggestions && users.length > 0 && (
                                <div className="user-suggestions">
                                    {users.map(user => (
                                        <div
                                            key={user.id}
                                            onClick={() => handleUserSelect(user)}
                                            className="user-suggestion-item"
                                        >
                                            <div className="user-avatar-small">
                                                {user.nickname.charAt(0).toUpperCase()}
                                            </div>
                                            <div className="user-info">
                                                <div className="user-name">{user.nickname}</div>
                                            </div>
                                        </div>
                                    ))}
                                </div>
                            )}
                        </div>
                    </div>

                    {/* Submit Button */}
                    <div className="form-actions">
                        <button
                            type="button"
                            className="cancel-group-btn"
                            onClick={onCancel}
                        >
                            Cancel
                        </button>
                        <button
                            type="submit"
                            className="submit-group-btn"
                            disabled={!groupTitle.trim()}
                        >
                            Create Group
                        </button>
                    </div>
                </form>
            </div>
        </div>
    );
}

export function GroupCreationTrigger() {
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [showNoGroups, setShowNoGroups] = useState(true);
    const [userList, setUserList] = useState([]);
    const [userID, setUserID] = useState('')

    const handlePostClick = () => {
        setShowNoGroups(false);
        setIsModalOpen(true);
    };
    const handleCancel = () => {
        setIsModalOpen(false);
        setShowNoGroups(true);
    };
    useEffect(() => {
        const fetchData = async () => {
            try {
                const resUser = await fetch('http://localhost:8080/api/me', {
                    credentials: 'include',
                });
                const userData = await resUser.json();
    
                const resFollowers = await fetch(`http://localhost:8080/api/followers?id=${userData.user_id}`, {
                    method: 'GET',
                    credentials: 'include',
                });
                const followersData = await resFollowers.json();
    
                setUserList(followersData);
            } catch (error) {
                console.error("Failed to fetch users when creating a group:", error);
            }
        };

        fetchData();
    }, []);
    

    console.log("usssssseeers are :", userList);
    return (
        <>
            <div className="create-post-container">
                <div className="create-post-card">
                    <div className="create-post-header">
                        <div className="user-avatar">
                            <span>User</span>
                        </div>
                        <input
                            type="text"
                            placeholder="DO you want to create a group ? just click here !"
                            className="post-input"
                            onClick={handlePostClick}
                            readOnly
                        />
                    </div>
                </div>
            </div>
            {isModalOpen && (
                <CreateGroupForm
                    users={userList}
                    onSubmit={createGroup}
                    onCancel={handleCancel}
                />
            )}
            {/* {showNoGroups && (
                <div className="no-groups">
                    You don’t have any groups yet. Click above to create one!
                </div>
            )} */}

        </>
    )
}

