import { useState, useEffect } from "react";
import "../styles/groupstyle.css";
import { createGroup } from "../app/groups/page";
import { GroupCard } from "./groupCard";
import { Plus } from "lucide-react";
import { Toaster, toast } from "sonner"
import { useWS } from "../context/wsContext";
export function CreateGroupForm({ users, onSubmit, onCancel }) {

  const [groupTitle, setGroupTitle] = useState("");
  const [groupDescription, setGroupDescription] = useState("");
  const [searchQuery, setSearchQuery] = useState("");
  const [selectedUsers, setSelectedUsers] = useState([]);
  const [showSuggestions, setShowSuggestions] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);


  const handleUserSelect = (user) => {
    if (!selectedUsers.find((u) => u.id === user.id)) {
      setSelectedUsers([...selectedUsers, user]);
      setSearchQuery("");
    }
  };

  const handleUserRemove = (userId) => {
    setSelectedUsers(selectedUsers.filter((u) => u.id !== userId));
  };

  const handleSubmit = (e) => {
    e.preventDefault();

    const groupData = {
      title: groupTitle,
      description: groupDescription,
      invitedUsers: selectedUsers.map((u) => u.id),
    };

    onSubmit(groupData, isSubmitting, setIsSubmitting);
    // onSubmit(selectedUsers);

    // Reset form after submission
    setGroupTitle("");
    setGroupDescription("");
    setSelectedUsers([]);
    setSearchQuery("");
  };


  return (
    <div className="group-modal-overlay">
      {/* <Toaster position="top-center" /> */}
      <div className="group-modal-content">
        <div className="group-modal-header">
          <h1 className="group-modal-title">Create a New Group</h1>
        </div>

        <form onSubmit={handleSubmit}>
          {/* Group Title */}
          <div className="group-modal-form">
            <div className="form-group">
              <label htmlFor="groupTitle" className="form-label">
                Group Title
              </label>
              <input
                id="groupTitle"
                type="text"
                value={groupTitle}
                onChange={(e) => setGroupTitle(e.target.value)}
                placeholder="Enter group title"
                className="form-input"
                required
                autoFocus
              />
            </div>

            {/* Group Description */}
            <div className="form-group">
              <label htmlFor="groupDescription" className="form-label">
                Description
              </label>
              <textarea
                id="groupDescription"
                value={groupDescription}
                onChange={(e) => setGroupDescription(e.target.value)}
                placeholder="What's your group about?"
                className="form-textarea"
                rows={4}
                required
              />
            </div>

            {/* Invite Users */}
            <fieldset className="form-group">
              <legend className="form-label">Invite Members</legend>

              {/* Selected Users */}
              {selectedUsers.length > 0 && (
                <ul className="selected-users" role="list">
                  {selectedUsers.map((user) => (
                    <li key={user.id} className="selected-user-tag">
                      <span className="user-avatar-small" aria-hidden="true">
                        {user.first_name.charAt(0).toUpperCase()}
                      </span>
                      <span>{user.first_name}</span>
                      <button
                        type="button"
                        onClick={() => handleUserRemove(user.id)}
                        className="remove-user-btn"
                        aria-label={`Remove ${user.name}`}
                      >
                        ✕
                      </button>
                    </li>
                  ))}
                </ul>
              )}

              {/* Search Input */}
              <div className="search-container">
                <input
                  id="inviteUsers"
                  type="text"
                  value={searchQuery}
                  onClick={() => setShowSuggestions((prev) => !prev)}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  placeholder="Search users to invite..."
                  className="form-input search-input"
                  autoComplete="off"
                />


                {/* User Suggestions */}
                {showSuggestions && users?.length > 0 && (
                  <ul
                    id="user-suggestions-list"
                    className="user-suggestions"
                    role="listbox"
                    aria-label="User suggestions"
                  >
                    {users.map((user) => (
                      <div
                        key={user.id}
                        id={`user-option-${user.id}`}
                        onClick={() => handleUserSelect(user)}
                        className="user-suggestion-item"
                      >
                        <span className="user-avatar-small" aria-hidden="true">
                          {user.first_name.charAt(0).toUpperCase()}
                        </span>
                        <div className="user-info">
                          <div className="user-name">{user.first_name} {user.last_name}</div>
                        </div>
                      </div>
                    ))}
                  </ul>
                )}
              </div>
            </fieldset>
          </div>

          {/* Submit Button */}
          <div className="group-modal-actions">
            <button type="button" className="cancel-button" onClick={onCancel}>
              Cancel
            </button>
            <button
              type="submit"
              className="submit-button"
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

export function GroupCreationTrigger({ setGroup }) {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [showNoGroups, setShowNoGroups] = useState(true);
  const [userList, setUserList] = useState([]);
  const { sendMessage } = useWS();



  const handlePostClick = () => {
    setShowNoGroups(false);
    setIsModalOpen(true);
  };
  const handleCancel = () => {
    setIsModalOpen(false);
    setShowNoGroups(true);
  };

  const handleSubmit = async (groupData, isSubmitting, setIsSubmitting) => {
    try {
      if (isSubmitting) return;
      setIsSubmitting(true);
      const newGroup = await createGroup(groupData, sendMessage);
      toast.success("group created successfully!");
      setGroup((prev) => {
        const exists = prev.some((g) => g.ID === newGroup.ID);
        return exists ? prev : [newGroup, ...prev];
      });
    } catch (error) {
      // console.error("Error creating group:", error);
      toast.error(error.message);
    } finally {
      setIsModalOpen(false);
      setIsSubmitting(false);
    }
  };
  useEffect(() => {
    const fetchData = async () => {
      try {
        const resUser = await fetch("http://localhost:8080/api/me", {
          credentials: "include",
        });
        const userData = await resUser.json();
        const resFollowers = await fetch(
          `http://localhost:8080/api/followers?id=${userData.user_id}`,
          {
            method: "GET",
            credentials: "include",
          }
        );
        const followersData = await resFollowers.json();
        setUserList(followersData);
      } catch (error) {
        console.error("Failed to fetch users when creating a group:", error);
      }
    };
    fetchData();
  }, []);

  return (
    <div className="create-card">
      <div className="create-card-inner">
        <div className="group-avatar">U</div>
        <input
          type="text"
          placeholder="Create a group?"
          className="create-input"
          readOnly
          onClick={handlePostClick}
        />
        <button onClick={handlePostClick} className="create-button">
          <Plus className="create-icon" />
        </button>
      </div>
      {isModalOpen && (
        <CreateGroupForm
          users={userList}
          onSubmit={handleSubmit}
          onCancel={handleCancel}
        />
      )}
    </div>
  );
}
