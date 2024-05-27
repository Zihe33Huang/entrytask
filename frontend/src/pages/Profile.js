import {
    Box,
    Modal,
    Slider,
    Button,
    Backdrop,
    CircularProgress,
} from "@mui/material";
import React, {useEffect, useRef, useState} from "react";
import AvatarEditor from "react-avatar-editor";
import "./Cropper.scss";
import {fetchWithJWT} from "../config";
import {NotificationManager} from "react-notifications";
import "./Profile.css"
import {FaEdit} from "react-icons/fa";
import {AiOutlineCheck, AiOutlineClose} from "react-icons/ai";
import notification from "react-notifications/lib/Notification";
import notificationManager from "react-notifications/lib/NotificationManager";
import {useNavigate} from "react-router-dom";


const style = {
    py: 0,
    width: '100%',
    maxWidth: 360,
    borderRadius: 2,
    border: '1px solid',
    borderColor: 'divider',
    backgroundColor: 'background.paper',
    marginTop: "50px",
};

// Styles
const boxStyle = {
    width: "300px",
    height: "300px",
    display: "flex",
    flexFlow: "column",
    justifyContent: "center",
    alignItems: "center"
};
const modalStyle = {
    display: "flex",
    justifyContent: "center",
    alignItems: "center"
};

// Modal
const CropperModal = ({ src, modalOpen, setModalOpen, setPreview }) => {
    const [slideValue, setSlideValue] = useState(10);
    const cropRef = useRef(null);

    //handle save
    const handleSave = async () => {
        if (cropRef) {
            const dataUrl = cropRef.current.getImage().toDataURL();
            const result = await fetch(dataUrl);
            const blob = await result.blob();
            setPreview(URL.createObjectURL(blob));
            setModalOpen(false);
        }
    };

    return (
        <Modal sx={modalStyle} open={modalOpen}>
            <Box sx={boxStyle}>
                <AvatarEditor
                    ref={cropRef}
                    image={src}
                    style={{ width: "100%", height: "100%" }}
                    border={50}
                    borderRadius={150}
                    color={[0, 0, 0, 0.72]}
                    scale={slideValue / 10}
                    rotate={0}
                />

                {/* MUI Slider */}
                <Slider
                    min={10}
                    max={50}
                    sx={{
                        margin: "0 auto",
                        width: "80%",
                        color: "cyan"
                    }}
                    size="medium"
                    defaultValue={slideValue}
                    value={slideValue}
                    onChange={(e) => setSlideValue(e.target.value)}
                />
                <Box
                    sx={{
                        display: "flex",
                        padding: "10px",
                        border: "3px solid white",
                        background: "black"
                    }}
                >
                    <Button
                        size="small"
                        sx={{ marginRight: "10px", color: "white", borderColor: "white" }}
                        variant="outlined"
                        onClick={(e) => setModalOpen(false)}
                    >
                        cancel
                    </Button>
                    <Button
                        sx={{ background: "#5596e6" }}
                        size="small"
                        variant="contained"
                        onClick={handleSave}
                    >
                        Save
                    </Button>
                </Box>
            </Box>
        </Modal>
    );
};

// Container
const Cropper = () => {
    // image src
    const [src, setSrc] = useState(null);

    // preview

    // modal state
    const [modalOpen, setModalOpen] = useState(false);

    // ref to control input element
    const inputRef = useRef(null);

    const [file, setFile] = useState(null);



    const [image, setImage] = useState(null);

    const [open, setOpen] = React.useState(false);

    const [isEditing, setIsEditing] = useState(false);

    const [updateNickname, setUpdateNickname] = React.useState("");

    const navigate = useNavigate();


    // useEffect(() => {
    //     setUpdateNickname(profileData.nickname)
    // }, [profileData]);


    // useEffect(() => {
    //     console.log('hhhh' + value)
    // }, [value]);
    const handleEditClick = () => {
        setUpdateNickname(profileData.nickname)
        setIsEditing(true);
    };


    const handleChange = (e) => {
        // setProfileData({...profileData, nickname: e.target.value});
        setUpdateNickname(e.target.value)
    };

    const handleNotSave = () => {
        console.log("cancel save")
        setUpdateNickname(profileData.nickname)
        setIsEditing(false)
    }

    // update nickname
    const handleSave = () => {
        handleOpen()
        console.log('save')
        const req = {
            nickname: updateNickname
        };
        fetchWithJWT('http://localhost:8080/api/users/nickname', {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(req)
        })
            .then(response => {
                console.log(response)
                if (response.ok) {
                    response.json().then(data => {
                        console.log(data)
                        if (data.status === 200) {
                           getProfile()
                            NotificationManager.success(data.message, "", 1000)
                        } else {
                            NotificationManager.error(data.message, "", 1000)
                        }
                        handleClose()
                        setIsEditing(false);
                    })
                } else {
                    // Handle login error
                    NotificationManager.warning("Server Error")
                    handleClose()
                    setIsEditing(false);
                }
            })
            .catch(error => {
                console.error('Error occurred during updating:', error);
            })
    };


    const [profileData, setProfileData] = useState({
        username: "",
        nickname: "",
    });


    const handleClose = () => {
        setOpen(false);
    };
    const handleOpen = () => {
        setOpen(true);
    };



    // handle Click
    const handleInputClick = (e) => {
        // setFile(e.target.files[0]);
        e.preventDefault();
        inputRef.current.click();
    };
    // handle Change
    const handleImgChange = (e) => {
        // setSrc(URL.createObjectURL(e.target.files[0]));
        // setModalOpen(true);
        setFile(e.target.files[0]);
    };

    useEffect(() => {
        if (file) {
            handleUpload();
        }
    }, [file]);

    useEffect(() => {
        getProfile()
    }, []);


    const getAvatar = () => {
        const url = `http://localhost:8080/api/users/avatar`;

        fetchWithJWT(url)
            .then(response => response.blob())
            .then(blob => {
                const url = URL.createObjectURL(blob);
                setImage(url);
            })
            .catch(error => console.error(error));
    };

    const getProfile = () => {
        fetchWithJWT('http://localhost:8080/api/users/profile')
            .then(response => response.json())
            .then(data => {
                if (data.status === 200) {
                    setProfileData(data.data);
                    getAvatar()
                } else {
                    NotificationManager.error("Please Log In!");
                    navigate("/")
                }
            })
            .catch(error => console.error('Error fetching profile data:', error));
    }


    const handleUpload = () => {
        if (!file) {
            console.log("Please select a file")
            return;
        }
        handleOpen()

        const formData = new FormData();
        formData.append('file', file);

        fetchWithJWT('http://localhost:8080/api/users/avatar', {
            method: 'POST',
            body: formData,
        })
            .then((response) => {
                return response.json();
            })
            .then((data) => {
                handleClose()
                if (data.status === 200) {
                    NotificationManager.success("Avatar updates successfully", "", 1000)
                    getAvatar()
                }
            })
            .catch((error) => {
                handleClose()
                NotificationManager.warning("Avatar update fails", "", 1000)
            });
    };


    return (
        <>
            <header>
                <div className="header-container">
                    <h1>Profile</h1>
                    <button onClick={() => {
                        localStorage.clear()
                        notificationManager.success("Log Out Successfully")
                        navigate("/")
                    }}>log out</button>
                </div>
                <hr />
            </header>
            <main className="container">
                <CropperModal
                    modalOpen={modalOpen}
                    src={src}
                    setModalOpen={setModalOpen}
                />
                <input
                    type="file"
                    accept="image/*"
                    ref={inputRef}
                    onChange={handleImgChange}
                />
                <div className="img-container">
                    <a onClick={handleInputClick}>
                        <img
                            src={
                                image
                            }
                            alt=""
                            width="200"
                            height="200"
                        ></img>
                    </a>
                </div>
                <div className="form">
                    <p>Username：<span>{profileData ? profileData.username : ""}</span></p>
                    <div className="nickname">
                        Nickname：
                        {isEditing ? (
                            <input
                                type="text"
                                value={updateNickname}
                                onChange={handleChange}
                                autoFocus
                            />
                        ) : (
                            <span>{profileData.nickname}</span>
                        )}
                        {!isEditing && (
                            <button onClick={handleEditClick}>
                                <FaEdit/>
                            </button>
                        )}
                        {isEditing && (
                            <div>
                                <button onClick={handleSave}>
                                    <AiOutlineCheck/>
                                </button>
                                <button onClick={handleNotSave}>
                                    <AiOutlineClose/>
                                </button>
                            </div>
                        )}
                    </div>
                </div>
            </main>
            <Backdrop
                sx={{color: '#fff', zIndex: (theme) => theme.zIndex.drawer + 1}}
                open={open}
            >
                <CircularProgress color="inherit"/>
            </Backdrop>

        </>
    );
};

export default Cropper;
