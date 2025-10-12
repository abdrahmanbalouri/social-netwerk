import React, { useRef, useState } from 'react';
import '../styles/ProfileCardEditor.css';

export default function ProfileCardEditor({
    showPrivacy,
    initialCover = '',
    initialAvatar = '',
    initialPrivacy = 'public',
    onSave = (state) => console.log('saved', state),

}) {
    const [coverPreview, setCoverPreview] = useState(initialCover);
    const [avatarPreview, setAvatarPreview] = useState(initialAvatar);
    const [privacy, setPrivacy] = useState(initialPrivacy);
    const [displayName, setDisplayName] = useState('about');
    const coverInputRef = useRef(null);
    const avatarInputRef = useRef(null);
    const [saving, setSaving] = useState(false);

    function handleFileToPreview(file, setPreview) {
        if (!file) return;
        const reader = new FileReader();
        reader.onload = (e) => setPreview(e.target.result);
        reader.readAsDataURL(file);
    }

    function onCoverChange(e) {
        const file = e.target.files?.[0];
        handleFileToPreview(file, setCoverPreview);
    }

    function onAvatarChange(e) {
        const file = e.target.files?.[0];
        handleFileToPreview(file, setAvatarPreview);
    }

    function clearCover() {
        setCoverPreview('');
        if (coverInputRef.current) coverInputRef.current.value = null;
    }

    function clearAvatar() {
        setAvatarPreview('');
        if (avatarInputRef.current) avatarInputRef.current.value = null;
    }

    async function handleSave() {
        setSaving(true);
        const formData = new FormData();
        formData.append("displayName", displayName);
        formData.append("privacy", privacy);

        if (coverInputRef.current?.files[0]) {
            formData.append("cover", coverInputRef.current.files[0]);
        }
        if (avatarInputRef.current?.files[0]) {
            formData.append("avatar", avatarInputRef.current.files[0]);
        }

        try {
            await fetch("http://localhost:8080/api/editor", {
                method: "POST",
                credentials: "include",
                body: formData,
            });

            onSave({ displayName, privacy });
            handleShowPrivacy()
        } finally {
            setSaving(false);
        }
    }
    function handleShowPrivacy() {
        
        showPrivacy = false

    }

    return (
        <div className="profile-card">
            {/* Cover area */}
            <div className="cover">
                {coverPreview ? (
                    <img src={coverPreview} alt="cover" />
                ) : (
                    <div className="placeholder">Cover Image</div>
                )}

                <div className="cover-actions">
                    <label className="btn">
                        Change
                        <input
                            ref={coverInputRef}
                            type="file"
                            accept="image/*"
                            className="hidden-input"
                            onChange={onCoverChange}
                        />
                    </label>
                    <button onClick={clearCover} type="button" className="btn">
                        Remove
                    </button>
                </div>
            </div>

            {/* Body */}
            <div className="body">
                <div className="avatar-section">
                    <div className="avatar-wrapper">
                        {avatarPreview ? (
                            <img src={avatarPreview} alt="avatar" className="avatar" />
                        ) : (
                            <div className="placeholder">Avatar</div>
                        )}

                        <div className="avatar-actions">
                            <label className="btn small">
                                Upload
                                <input
                                    ref={avatarInputRef}
                                    type="file"
                                    accept="image/*"
                                    className="hidden-input"
                                    onChange={onAvatarChange}
                                />
                            </label>
                            <button onClick={clearAvatar} className="btn small">
                                Remove
                            </button>
                        </div>
                    </div>

                    <div className="name-section">
                        <input
                            value={displayName}
                            onChange={(e) => setDisplayName(e.target.value)}
                            className="display-name"
                            aria-label="Display Name"
                        />
                        <p className="subtext">Short bio or user title</p>
                    </div>
                </div>

                <div className="settings">
                    <label className="field">
                        <span>Profile Privacy</span>
                        <select
                            value={privacy}
                            onChange={(e) => setPrivacy(e.target.value)}
                            className="select"
                        >
                            <option value="public">Public</option>
                            <option value="friends">Friends</option>
                            <option value="private">Private</option>
                        </select>
                    </label>

                    

                    <div className="actions">
                        <button
                            onClick={() => {
                                setCoverPreview(initialCover);
                                setAvatarPreview(initialAvatar);
                                setPrivacy(initialPrivacy);
                            }}
                            className="btn"
                        >
                            Cancel
                        </button>

                        <button onClick={handleSave} className="btn primary" disabled={saving}>
                            {saving ? 'Saving...' : 'Save changes'}

                        </button>
                    </div>
                </div>
            </div>
        </div>
    );
}
