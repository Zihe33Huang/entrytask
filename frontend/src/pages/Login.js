import React, { useState } from 'react';

import {
    Container,
    CssBaseline,
    Avatar,
    Typography,
    TextField,
    Button,
    Grid,
    Link,
    Paper,
    IconButton,
    InputAdornment,
    Backdrop, CircularProgress
} from '@mui/material';
import VisibilityIcon from '@mui/icons-material/Visibility';
import VisibilityOffIcon from '@mui/icons-material/VisibilityOff';
import { styled } from '@mui/system';
import 'react-notifications/lib/notifications.css';
import {NotificationContainer, NotificationManager} from "react-notifications";
import {useNavigate} from "react-router-dom";


const CenteredContainer = styled(Container)({
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    justifyContent: 'center',
    minHeight: '100vh',
});

const AnimatedPaper = styled(Paper)({
    padding: '2rem',
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    gap: '1rem',
    transition: 'transform 0.5s ease-in-out',
    '&:hover': {
        transform: 'scale(1.1)',
    },
});

const LoginPage = () => {
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [showPassword, setShowPassword] = useState(false);
    const [error, setError] = useState('');
    const [open, setOpen] = React.useState(false);
    const navigate = useNavigate();

    const handleKeyPress = (event) => {
        if (event.key === 'Enter') {
            handleLogin()
        }
    };


    const handleClose = () => {
        setOpen(false);
    };
    const handleOpen = () => {
        setOpen(true);
    };

    const handleLogin = () => {
        const loginData = {
            username: username,
            password: password
        };
        handleOpen()
        fetch('http://localhost:8080/api/users/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(loginData)
        })
            .then(response => {
                console.log(response)
                if (response.ok) {
                    response.json().then(data => {
                        console.log(data)
                        if (data.status === 200) {
                            const jwt =  response.headers.get('Authorization')
                            localStorage.clear()
                            localStorage.setItem('jwt', jwt);
                            // Login successful, navigate to the main page
                            NotificationManager.success(data.message, "", 1000)
                            navigate("/profile")
                        } else {
                            NotificationManager.error(data.message, "", 1000)
                        }
                        handleClose()
                    })
                } else {
                    // Handle login error
                    NotificationManager.warning("Server Error")
                    handleClose()
                }
            })
            .catch(error => {
                console.error('Error occurred during login:', error);
            });
    };

    const handleUserChange = (e) => {
        setUsername(e.target.value);
        setError('');
    };

    const handlePasswordChange = (e) => {
        setPassword(e.target.value);
        setError('');
    };

    const togglePasswordVisibility = () => {
        setShowPassword(!showPassword);
    };

    return (
        <div onKeyUp={handleKeyPress}>
            <CenteredContainer component="main" maxWidth="xs">
                <CssBaseline />
                <AnimatedPaper elevation={5}>
                    <Typography component="h1" variant="h5">
                        Log in
                    </Typography>
                    <form style={{ width: '100%' }}>
                        <Grid container spacing={2}>
                            <Grid item xs={12}>
                                <TextField
                                    fullWidth
                                    label="Username"
                                    variant="outlined"
                                    value={username}
                                    onChange={handleUserChange}
                                />
                            </Grid>
                            <Grid item xs={12}>
                                <TextField
                                    fullWidth
                                    label="Password"
                                    type={'text'}
                                    variant="outlined"
                                    value={password}
                                    onChange={handlePasswordChange}
                                    InputProps={{
                                        endAdornment: (
                                            <InputAdornment position="end">
                                                <IconButton onClick={togglePasswordVisibility} edge="end">
                                                    {/*{showPassword ? <VisibilityIcon /> : <VisibilityOffIcon />}*/}
                                                </IconButton>
                                            </InputAdornment>
                                        ),
                                    }}
                                />
                            </Grid>
                        </Grid>
                        {error && (
                            <Typography variant="body2" color="error" style={{ marginTop: '0.5rem' }}>
                                {error}
                            </Typography>
                        )}
                        <Button
                            type="button"
                            fullWidth
                            variant="contained"
                            style={{ backgroundColor: error ? '#ffffff' : '#4CAF50', color: 'black', marginTop: '1rem' }}
                            onClick={handleLogin}
                            disabled={error !== ''}
                        >
                            Log In
                        </Button>
                    </form>
                </AnimatedPaper>
            </CenteredContainer>
            <Backdrop
                sx={{ color: '#fff', zIndex: (theme) => theme.zIndex.drawer + 1 }}
                open={open}
            >
                <CircularProgress color="inherit" />
            </Backdrop>
        </div>
    );
};

export default LoginPage;