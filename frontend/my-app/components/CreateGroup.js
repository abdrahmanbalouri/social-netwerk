import { useState, useEffect } from 'react';
import "../styles/createGroup.css"

export function CreateGroupForm({ users, onSubmit, onCancel }) {
    const [groupTitle, setGroupTitle] = useState('');
    const [groupDescription, setGroupDescription] = useState('');
    const [searchQuery, setSearchQuery] = useState('');
    const [selectedUsers, setSelectedUsers] = useState([]);
    const [filteredUsers, setFilteredUsers] = useState([]);

    // Filter users based on search query
    useEffect(() => {
        if (searchQuery.trim()) {
            const filtered = users.filter(user =>
                user.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
                user.username.toLowerCase().includes(searchQuery.toLowerCase())
            );
            setFilteredUsers(filtered);
        } else {
            setFilteredUsers([]);
        }
    }, [searchQuery, users]);

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

        // Reset form after submission
        setGroupTitle('');
        setGroupDescription('');
        setSelectedUsers([]);
        setSearchQuery('');
    };

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
                                            {user.name.charAt(0).toUpperCase()}
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
                                onChange={(e) => setSearchQuery(e.target.value)}
                                placeholder="Search users to invite..."
                                className="form-input search-input"
                            />

                            {/* User Suggestions */}
                            {filteredUsers.length > 0 && (
                                <div className="user-suggestions">
                                    {filteredUsers.map(user => (
                                        <div
                                            key={user.id}
                                            onClick={() => handleUserSelect(user)}
                                            className="user-suggestion-item"
                                        >
                                            <div className="user-avatar-small">
                                                {user.name.charAt(0).toUpperCase()}
                                            </div>
                                            <div className="user-info">
                                                <div className="user-name">{user.name}</div>
                                                <div className="user-username">@{user.username}</div>
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

    const handlePostClick = () => {
        setShowNoGroups(false);
        setIsModalOpen(true);
    };
    const handleCancel = () => {
        setIsModalOpen(false);
        setShowNoGroups(true);
    };

    const mockSubmitHandler = (formData) => {
        console.log('--- Mock Submission Attempt ---');
        console.log('✅ Form successfully collected and submitted the following data:');

        // Simulate a network delay for better testing experience
        setTimeout(() => {
            console.log(formData);
            console.log('--- Submission Mock Finished ---');
            alert(`Form Submitted! Check the browser console for data.\nName: ${formData.groupName}`);
        }, 500);
    };
    const userList = []
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
                    onSubmit={mockSubmitHandler}
                    onCancel={handleCancel}
                />
            )}
            {showNoGroups && (
                <div className="no-groups">
                    You don’t have any groups yet. Click above to create one!
                </div>
            )}

        </>
    )
}

