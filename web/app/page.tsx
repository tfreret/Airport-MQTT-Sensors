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

interface DataMeans {
  TempAverage: string,
  PresAverage: string,
  WindAverage: string,
}

export default function Home() {
  const [airports, setAirports] = useState<string[]>([]);
  const [sensors, setSensors] = useState<Sensor[]>([]);
  const [selectedAirport, setSelectedAirport] = useState("");
  const [selectedSensor, setSelectedSensor] = useState("");
  const [selectedMeans, setSelectedMeans] = useState<DataMeans>();
  const [activeTab, setActiveTab] = useState(0);

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

  const fetchAverage = useMemo(async () => {
    if (selectedAirport === "") return
    try {
      const response = await axios.get(`http://localhost:8080/average/${selectedAirport}`);
      console.log(response.data)
      setSelectedMeans(response.data);
    } catch (err) {
      console.log(err);
    }
  }, [selectedAirport]);

  const handleSensorChange = useCallback((e: string) => {setSelectedSensor(e)}, []);
  const handleTabChange = (index: number) => {setActiveTab(index);};

  return (
    <main className="p-12">
      <Title>Tableau de bord</Title>
      <Card className="mt-6">
      <Flex justifyContent="start" className="gap-10">
        <Title>Aéroport :</Title>
        <SearchSelect value={selectedAirport} onValueChange={setSelectedAirport} className="w-auto">
            {airports.map(it=><SearchSelectItem value={it} key={it}></SearchSelectItem>)}
        </SearchSelect>
        {activeTab === 1 && (
          <>
            <Title>Capteur :</Title>
            <SearchSelect value={selectedSensor} onValueChange={handleSensorChange} disabled={selectedAirport === ""} className="w-auto">
            {sensors.map(it=><SearchSelectItem value={it.ID} key={it.ID}></SearchSelectItem>)}
            </SearchSelect>
          </>
        )}

      </Flex>
      </Card>
      <TabGroup className="mt-6">
        <TabList>
          <Tab onClick={() => handleTabChange(0)}>Aéroport</Tab>
          <Tab onClick={() => handleTabChange(1)}>Capteur</Tab>
        </TabList>
        <TabPanels>
          <TabPanel>
            {selectedAirport ? (
            <Grid numItemsMd={2} numItemsLg={3} className="gap-6 mt-6">
              <Card>
                <Flex justifyContent="center" alignItems="center" flexDirection="col" className="gap-5" >
                <Title>Nombre de capteurs :</Title>
                <Metric>{sensors.length}</Metric>
                </Flex>
              </Card>
              <Card>
                <Flex justifyContent="center" alignItems="center" flexDirection="col" className="gap-5" >
                <Title>Température moyenne (aujourd'hui) :</Title>
                <Metric>{selectedMeans?.TempAverage + " °C"}</Metric>
                </Flex>
              </Card>
              <Card>
              <Flex justifyContent="center" alignItems="center" flexDirection="col" className="gap-5" >
                <Title>Vitesse moyenne du vent (aujourd'hui) :</Title>
                <Metric>{selectedMeans?.WindAverage + " km/h"}</Metric>
                </Flex>
              </Card>
              <Card>
                <Flex justifyContent="center" alignItems="center" flexDirection="col" className="gap-5" >
                <Title>Pression atmosphérique moyenne (aujourd'hui) :</Title>
                <Metric>{selectedMeans?.PresAverage + " hPa"}</Metric>
                </Flex>
              </Card>
            </Grid>
            ) : (
              <Grid numItemsMd={2} numItemsLg={3} className="gap-6 mt-6">
              <Card>
                <Flex justifyContent="center" alignItems="center" flexDirection="col" className="gap-5" >
                <Title>Aucun aéroport sélectionné</Title>
                </Flex>
              </Card>
              </Grid>
            )}

          </TabPanel>
          <TabPanel>
          <div className="mt-6">
              { selectedAirport !== "" && selectedSensor !== "" ? (
                <Chart key={`${selectedAirport}/${selectedSensor}`} url={`http://localhost:8080/data/${selectedAirport}/${sensors.find(it=>it.ID===selectedSensor)?.MeasureType}/${selectedSensor}`} />
              ) : (
                <Grid numItemsMd={2} numItemsLg={3} className="gap-6 mt-6">
                  <Card>
                    <Flex justifyContent="center" alignItems="center" flexDirection="col" className="gap-5" >
                      <Title>Aucun aéroport/capteur sélectionné</Title>
                    </Flex>
                  </Card>
                </Grid>
              )
              }
            </div>
          </TabPanel>
        </TabPanels>
      </TabGroup>
    </main>
  );
}