import React, { useState } from "react";
import { Box, Button, Center, Container, FormControl, HStack, Image, Input } from "@chakra-ui/react";
import Navbar from "../../components/Navbar/Navbar";
import "./home.css"

const getRequestOptions = (img) => {
    const raw = JSON.stringify({
        "user_app_id": {
            "user_id": "neeejm",
            "app_id": "6914fa6e4c2544f982cb13b63d9874f6"
        },
        "inputs": [
            {
                "data": {
                    "image": {
                        "url": `${img}`
                    }
                }
            }
        ]
    });

    const requestOptions = {
        method: 'POST',
        headers: {
            'Accept': 'application/json',
            'Authorization': 'Key eae5e9c33026413aadca8088320b7880'
        },
        body: raw
    };

    return requestOptions;
}

const Home = () => {
    const [imgSrc, setImgSrc] = useState("");
    const [boundingBox, setBoundingBox] = useState({});

    const img = (() => {
        if (imgSrc !== "") {
            return (
                <Center>
                    <Image id="face" src={imgSrc} maxW="550px" maxH="550px" paddingBottom="30px" />
                    <div className="bounding-box"
                        style={{ top: boundingBox.top, right: boundingBox.right, bottom: boundingBox.bottom, left: boundingBox.left }}>
                    </div>
                </Center>
            );

        }
    })();

    window.addEventListener("resize", () => {
        detectFace();
    });

    const handleImg = (event) => {
        try {
            new URL(event.target.value);
            setBoundingBox({});
            setImgSrc(event.target.value);
        }
        catch (err) {
            console.log(err);
        }
    }

    const detectFace = () => {
        fetch("https://api.clarifai.com/v2/models/f76196b43bbd45c99b4f3cd8e8b40a8a/versions/45fb9a671625463fa646c3523a3087d5/outputs", getRequestOptions(imgSrc))
            .then(response => response.json())
            // .then(result => console.log(result.outputs[0].data.regions[0].region_info.bounding_box))
            .then(result => calculateFaceLocation(result.outputs[0].data.regions[0].region_info.bounding_box))
            .catch(error => console.log('error', error));
    }

    const calculateFaceLocation = (faceLocation) => {
        const imgElement = document.getElementById("face");
        const width = Number(imgElement.width);
        const height = Number(imgElement.height);
        const offsetLeft = Number(imgElement.offsetLeft);
        const offsetTop = Number(imgElement.offsetTop);
        const vpHeight = window.innerHeight;

        setBoundingBox({
            top: offsetTop + (faceLocation.top_row * height), // check
            right: (width + offsetLeft) - (faceLocation.right_col * width), // check
            bottom: vpHeight - (vpHeight - offsetTop - faceLocation.bottom_row),
            left: offsetLeft + (faceLocation.left_col * width) // check
        });
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