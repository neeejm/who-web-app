import React, { useState } from "react";
import { Box, Button, Center, Container, FormControl, FormHelperText, FormLabel, Input, Text } from '@chakra-ui/react'
import { Link } from "react-router-dom";
import Navbar from "../../components/Navbar/Navbar";

const Auth = () => {

    const params = new URLSearchParams(window.location.search);
    const [reg, setReg] = useState(params.get("reg"));
    console.log("test ", window.location.href)
    // debugger;

    if (reg) {
        if ((reg.toLowerCase() === 'true') ? true : false) {
            console.log("signup page");
            return (
                <div>
                    <Navbar signIn={() => setReg("false")} signUp={() => setReg("true")} />
                    <Center>
                        <Box w="50%" maxW='90%' borderWidth='1px' borderRadius='lg'>
                            <Container maxW="container.sm" padding="30px">
                                <Text fontSize="2xl" color="blue">Sign Up Page</Text>

                                <FormControl>
                                    <FormLabel htmlFor='username'>Username</FormLabel>
                                    <Input id='username' type='text' />
                                    <FormHelperText>A validation mail will be sent.</FormHelperText>

                                    <FormLabel htmlFor='email'>Email address</FormLabel>
                                    <Input id='email' type='email' />
                                    <FormHelperText>A validation mail will be sent.</FormHelperText>

                                    <FormLabel htmlFor='password'>Password</FormLabel>
                                    <Input id='password' type='password' />
                                    <FormHelperText>Password must contain letters, numbers and symbols.</FormHelperText>

                                    <Button colorScheme="yellow">
                                        <Link to="/home">Sign Up</Link>
                                    </Button>
                                </FormControl>
                            </Container>
                        </Box>
                    </Center>
                </div>
            );
        }
    }

    console.log("signin page");
    return (
        < div >
            <Navbar signIn={() => setReg("false")} signUp={() => setReg("true")} />
            <Center>
                <Box w="50%" maxW='90%' borderWidth='1px' borderRadius='lg'>
                    <Container maxW="container.sm" padding="30px">
                        <Text fontSize="2xl" color="blue">Sign In Page</Text>

                        <FormControl>
                            <FormLabel htmlFor='email'>Email address</FormLabel>
                            <Input id='email' type='email' />

                            <FormLabel htmlFor='password'>Password</FormLabel>
                            <Input id='password' type='password' />

                            <Button colorScheme="yellow">
                                <Link to="/home">Sign In</Link>
                            </Button>
                        </FormControl>
                    </Container>
                </Box>
            </Center>
        </div >
    );
}

export default Auth;