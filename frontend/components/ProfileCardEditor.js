import { useRef, useState } from "react";
import "../styles/ProfileCardEditor.css";

export default function ProfileCardEditor({
    handleShowPrivacy,
    initialCover = "",
    initialAvatar = "",
    initialAbout = "about",
    initialPrivacy = "public",
    onSave = (state) => console.log("saved", state),
}) {
    const [coverPreview, setCoverPreview] = useState("");
    const [avatarPreview, setAvatarPreview] = useState("");
    const [privacy, setPrivacy] = useState(initialPrivacy);
    const [displayName, setDisplayName] = useState(initialAbout);
    const [toast, setToast] = useState(null)
    const coverInputRef = useRef(null);
    const avatarInputRef = useRef(null);
    const [saving, setSaving] = useState(false);
    const showToast = (message, type = "error", duration = 3000) => {
        setToast({ message, type });
        setTimeout(() => {
            setToast(null);
        }, duration);
    };
    function handleFileToPreview(file, setPreview) {
        if (!file) return;
        const allowedTypes = ["image/png", "image/jpeg", "image/jpg", "image/gif", "image/webp"];
        if (!allowedTypes.includes(file.type)) {
            showToast("Only image files are allowed!")
            e.target.value = "";
            setPreview(null);
            return;
        }
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
        setCoverPreview("");
        if (coverInputRef.current) coverInputRef.current.value = null;
    }

    function clearAvatar() {
        setAvatarPreview("");
        if (avatarInputRef.current) avatarInputRef.current.value = null;
    }

    async function handleSave() {
        setSaving(true);
        const formData = new FormData();

        formData.append("displayName", displayName || initialAbout);
        formData.append("privacy", privacy || initialPrivacy);

        if (coverInputRef.current.files[0]) {
            formData.append("cover", coverInputRef.current.files[0]);


        } else if (initialCover) {

            formData.append("existingCover", initialCover);
        }

        if (avatarInputRef.current?.files[0]) {
            formData.append("avatar", avatarInputRef.current.files[0]);
        } else if (initialAvatar) {
            formData.append("existingAvatar", initialAvatar);
        }

        try {
            await fetch("http://localhost:8080/api/editor", {
                method: "POST",
                credentials: "include",
                body: formData,
            });

            onSave(formData);
            handleShowPrivacy();
        } finally {
            setSaving(false);
        }
    }


    return (
        <div className="profile-card">
            {toast && (
                <div className={`toast ${toast.type}`}>
                    <span>{toast.message}</span>
                    <button onClick={() => setToast(null)} className="toast-close">×</button>
                </div>
            )}
            {/* Cover area */}
            <div className="cover">
                <img
                    src={
                        coverPreview
                            ? coverPreview
                            : initialCover
                                ? `/uploads/${initialCover}`
                                : "https://images.pexels.com/photos/13440765/pexels-photo-13440765.jpeg"
                    }
                    alt="cover"
                />

                <div className="cover-actions">
                    <label className="btn">
                        <i className="fa-solid fa-pen-to-square"></i>
                        <input
                            ref={coverInputRef}
                            type="file"
                            accept="image/*"
                            className="hidden-input"
                            onChange={onCoverChange}
                        />
                    </label>
                    <button onClick={clearCover} type="button" className="btn">
                        <i className="fa-solid fa-trash"></i>
                    </button>
                </div>
            </div>

            {/* Body */}
            <div className="body">
                <div className="avatar-section">
                    <div className="avatar-wrapper">
                        <img
                            src={
                                avatarPreview
                                    ? avatarPreview
                                    : initialAvatar
                                        ? `/uploads/${initialAvatar}`
                                        : "/assets/default.png"
                            }
                            className="avatar"
                        />

                        <div className="avatar-actions">
                            <label className="btn small">
                                <i className="fa-solid fa-pen-to-square"></i>
                                <input
                                    ref={avatarInputRef}
                                    type="file"
                                    accept="image/*"
                                    className="hidden-input"
                                    onChange={onAvatarChange}
                                />
                            </label>
                            <button onClick={clearAvatar} className="btn small">
                                <i className="fa-solid fa-trash"></i>
                            </button>
                        </div>
                    </div>

                    <div className="name-section">
                        <input
                            value={displayName}
                            onChange={(e) => setDisplayName(e.target.value)}
                            className="display-name"
                            aria-label="About"
                            id="rr"
                        />
                        <p className="subtext">Update your Short bio <i className="fa-solid fa-arrow-up"></i></p>
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
                            <option value="private">Private</option>
                        </select>
                    </label>

                    <div className="actions">
                        <button
                            onClick={handleSave}
                            className="btn primary"
                            disabled={saving}
                        >
                            {saving ? "Saving..." : "Save changes"}
                        </button>
                    </div>
                </div>
            </div>
        </div>
    );
}
