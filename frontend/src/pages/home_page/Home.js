import React, { useState } from "react";
import { Box, Button, Center, Container, FormControl, HStack, Image, Input, VStack, Text } from "@chakra-ui/react";
import Navbar from "../../components/Navbar/Navbar";
import "./home.css"

const IMAGE_PROXY_URL = "https://who-app.neeejm.workers.dev/?"
const ENCODE = "data:image/png;base64, "
let isBuffer = false

const Home = () => {
    const [imgSrc, setImgSrc] = useState("https://t3.ftcdn.net/jpg/02/57/43/20/360_F_257432094_IBWsGRGo9DXMS9glqquVlp3QQOly2UZA.jpg");
    const [faces, setFaces] = useState([]);
    const [isLoading, setIsLoading] = useState(null);

    const img = (() => {
        if (imgSrc !== "") {
            console.log("buffer: ", isBuffer)
            if (isBuffer) {
                if (isLoading)
                    return (
                        <Center>
                            <div className="loader"></div>
                        </Center>
                    );
                else
                    return (
                        <div>
                            <Center>
                                <Box>
                                    <Image id="face" src={ENCODE + imgSrc} maxW="550px" maxH="550px" />
                                    <div className="bounding-box">
                                    </div>
                                </Box>
                            </Center>
                            <Box padding="20px">
                                <VStack maxH="100px" overflowX="auto">
                                    {faces.map((face, index) => {
                                        return (
                                            <div>
                                                <HStack>
                                                    <Image id={index} src={ENCODE + face.image} maxW="100px" maxH="100px" />
                                                    <Text fontsize="md">
                                                        {face.gender ? ((face.name === "") ? "John Doe" : face.name) + " (male)" :
                                                            ((face.name === "") ? "Jane Doe" : face.name) + " (female)"}
                                                    </Text>
                                                </HStack>
                                            </div>
                                        );
                                    })}
                                </VStack>
                            </Box>
                        </div>
                    );
            }
            else {
                return (
                    <Center>
                        <Image id="face" src={IMAGE_PROXY_URL + imgSrc} maxW="550px" maxH="550px" />
                        <div className="bounding-box">
                        </div>
                    </Center>
                );
            }

        }
    })();

    const handleImg = (event) => {
        if (event.target.value === "") {
            setImgSrc("https://t3.ftcdn.net/jpg/02/57/43/20/360_F_257432094_IBWsGRGo9DXMS9glqquVlp3QQOly2UZA.jpg")
            setFaces([])
        }
        else {
            setImgSrc(event.target.value);
            isBuffer = false
        }
    }

    const detectFace = () => {
        try {
            new URL(imgSrc);
            setIsLoading(true);
            isBuffer = true;
            fetch("http://127.0.0.1:3003/detect?url=" + imgSrc, {
                // mode: "no-cors",
                method: 'POST'
            })
                .then(response => response.json())
                // .then(result => console.log(result.outputs[0].data.regions[0].region_info.bounding_box))
                .then(result => {
                    setIsLoading(false);
                    setImgSrc(result.images.image_box);
                    setFaces(result.images.faces)
                })
                // .then(result => console.log(result))
                .catch(error => console.log('error', error));
        }
        catch (err) {
            console.error("not a url");
            return;
        }
    }

    return (
        <div>
            <Navbar />
            <Container paddingTop="50px" maxW="container.sm">
                <Box borderWidth='4px' borderRadius='lg' overflow='hidden' padding="30px" marginBottom="20px" borderColor="#2e479c">
                    {img}
                </Box>
                <FormControl>
                    <Box >
                        <div style={{ paddingBottom: "100px" }}>
                            <HStack>
                                <Input placeholder='image of the face to find...' onChange={handleImg} borderColor="#2e479c" borderWidth="2px" />
                                <Button textTransform="uppercase" colorScheme="yellow" onClick={detectFace}>detect</Button>
                                {/* <Button colorScheme="red">find</Button> */}
                            </HStack>
                        </div>
                    </Box>
                </FormControl>
            </Container>
        </div >
    );
}

export default Home;