import React from "react";
import { Link } from "react-router-dom";
import { Flex, Image, Button, Spacer, VStack } from '@chakra-ui/react'

const Navbar = ({ signIn, signUp }) => {
    return (
        <Flex padding="35px">
            <Link to='/'>
                <Image
                    src="https://icon-library.com/images/facial-recognition-icon/facial-recognition-icon-26.jpg"
                    alt="a face getting recognized"
                    boxSize="128px"
                />
            </Link>
            <Spacer />
            <VStack>
                <Button colorScheme="blue" onClick={signUp}>
                    <Link to="/Auth?reg=true">Sign Up</Link>
                </Button>
                <Button colorScheme="gray" onClick={signIn}>
                    <Link to="/Auth">Sign In</Link>
                </Button>
            </VStack>
        </Flex>
    );
}

export default Navbar;