"use client";

      import {
        Card,
        Col,
        Flex,
        Grid,
        Metric,
        SearchSelect,
        SearchSelectItem,
        Subtitle,
        Tab,
        TabGroup,
        TabList,
        TabPanel,
        TabPanels,
        Text,
        Title,
      } from "@tremor/react";
import Chart from "./chart";
import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import axios from 'axios';

interface Sensor {
  ID: string,
  MeasureType: string
}

export default function Home() {
  const [airports, setAirports] = useState<string[]>([]);
  const [sensors, setSensors] = useState<Sensor[]>([]);
  const [selectedAirport, setSelectedAirport] = useState("");
  const [selectedSensor, setSelectedSensor] = useState("");

  const fetchAirports = useMemo(async () => {
      try {
        const response = await axios.get("http://localhost:8080/airports");
        setAirports(response.data);
      } catch (err) {
        console.log(err);
      }
  }, []);

  const fetchSensors = useMemo(async () => {
    if (selectedAirport === "") return
    try {
      const response = await axios.get(`http://localhost:8080/sensors/${selectedAirport}`);
      setSensors(response.data);
    } catch (err) {
      console.log(err);
    }
  }, [selectedAirport]);

  const handleSensorChange = useCallback((e: string) => {setSelectedSensor(e)}, []);

  return (
    <main className="p-12">
      <Title>Tableau de bord</Title>
      <Card className="mt-6">
      <Flex justifyContent="start" className="gap-10">
        <Title>Aéroport :</Title>
        <SearchSelect value={selectedAirport} onValueChange={setSelectedAirport} className="w-auto">
            {airports.map(it=><SearchSelectItem value={it} key={it}></SearchSelectItem>)}
          </SearchSelect>
          <Title>Capteur :</Title>
        <SearchSelect value={selectedSensor} onValueChange={handleSensorChange} disabled={selectedAirport === ""} className="w-auto">
          {sensors.map(it=><SearchSelectItem value={it.ID} key={it.ID}></SearchSelectItem>)}
        </SearchSelect>
      </Flex>
      </Card>
      <TabGroup className="mt-6">
        <TabList>
          <Tab>Aéroport</Tab>
          <Tab>Capteur</Tab>
        </TabList>
        <TabPanels>
          <TabPanel>
            <Grid numItemsMd={2} numItemsLg={3} className="gap-6 mt-6">
              <Card>
                <Flex justifyContent="center" alignItems="center" flexDirection="col" className="gap-5" >
                <Title>Nombre de capteurs :</Title>
                <Metric>{sensors.length}</Metric>
                </Flex>
                
              </Card>
              <Card>
                {/* Placeholder to set height */}
                <div className="h-28" />
              </Card>
              <Card>
                {/* Placeholder to set height */}
                <div className="h-28" />
              </Card>
            </Grid>
      
          </TabPanel>
          <TabPanel>
          <div className="mt-6">
              { selectedAirport !== "" && selectedSensor !== "" &&
              <Chart key={`${selectedAirport}/${selectedSensor}`} url={`http://localhost:8080/${sensors.find(it=>it.ID===selectedSensor)?.MeasureType}/${selectedAirport}/${selectedSensor}`} />
              }
            </div>
          </TabPanel>
        </TabPanels>
      </TabGroup>
    </main>
  );
}