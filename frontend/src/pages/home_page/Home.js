import React, { useState } from "react";
import { Box, Button, Center, Container, FormControl, HStack, Image, Input } from "@chakra-ui/react";
import Navbar from "../../components/Navbar/Navbar";
import "./home.css"

const IMAGE_PROXY_URL = "https://who-app.neeejm.workers.dev/?"
const ENCODE = "data:image/png;base64, "
let isBuffer = false

const Home = () => {
    const [imgSrc, setImgSrc] = useState("");

    const img = (() => {
        if (imgSrc !== "") {
            console.log("buffer: ", isBuffer)
            if (isBuffer) {
                return (
                    <Center>
                        <Image id="face" src={ENCODE + imgSrc} maxW="550px" maxH="550px" paddingBottom="30px" />
                        <div className="bounding-box">
                        </div>
                    </Center>
                );
            }
            else {
                return (
                    <Center>
                        <Image id="face" src={IMAGE_PROXY_URL + imgSrc} maxW="550px" maxH="550px" paddingBottom="30px" />
                        <div className="bounding-box">
                        </div>
                    </Center>
                );
            }

        }
    })();

    const handleImg = (event) => {
        try {
            // new URL(event.target.value);
            setImgSrc(event.target.value);
            isBuffer = false
        }
        catch (err) {
            console.log(err);
        }
    }

    const detectFace = () => {
        isBuffer = true
        fetch("http://127.0.0.1:3003/drawbox?url=" + imgSrc, {
            // mode: "no-cors",
            method: 'POST'
        })
            .then(response => response.json())
            // .then(result => console.log(result.outputs[0].data.regions[0].region_info.bounding_box))
            .then(result => setImgSrc(result.image.url))
            // .then(result => console.log(result))
            .catch(error => console.log('error', error));
    }

    return (
        <div>
            <Navbar />
            <Container paddingTop="50px" maxW="container.sm">
                {img}
                <FormControl>
                    <Box >
                        <div style={{ paddingBottom: "100px" }}>
                            <HStack>
                                <Input placeholder='image of the face to find...' onChange={handleImg} />
                                <Button colorScheme="yellow" onClick={detectFace}>detect</Button>
                                <Button colorScheme="red">find</Button>
                            </HStack>
                        </div>
                    </Box>
                </FormControl>
            </Container>
        </div >
    );
}

export default Home;