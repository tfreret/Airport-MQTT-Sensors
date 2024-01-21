import { AreaChart, Color, DateRangePicker, Flex, Icon, LineChart, Tab, TabGroup, TabList, Text, Title } from "@tremor/react";
import { useState, useEffect, useMemo } from "react";
import axios from 'axios';

interface Point {
  date: string
  value: number
}

interface PointAPI {
  Time: string,
  Value: number
}

import { InformationCircleIcon } from "@heroicons/react/solid";

interface Props {
  url: string
}
export default function Chart({url}: Props) {

  const [chartData, setChartData] = useState<Point[]>([]);

  const fetchData = async () => {
    try {
        const response = await axios.get(url);
        const data = response.data['tab of points'] as PointAPI[];
        setChartData(
          data.map(it=> {return {date: it.Time, value: it.Value}})
        );

    } catch (error) {
      setChartData([]);
      console.error('Erreur lors de la récupération des données:', error);
    }
  };

  useEffect(() => {
    fetchData();
    const intervalId = setInterval(() => {
        fetchData()
      }, 1000 * 5) // in milliseconds
      return () => clearInterval(intervalId)
  }, []);

  return (
    <>
      <div className="md:flex justify-between">
        <div>
          <Flex className="space-x-0.5" justifyContent="start" alignItems="center">
            <Title> Données du capteur </Title>
            <Icon
              icon={InformationCircleIcon}
              variant="simple"
              tooltip="Afficher les données du capteur"
            />
          </Flex>
        </div>
        <div>
          <DateRangePicker></DateRangePicker>
        </div>
      </div>
      <div className="mt-8 hidden sm:block">
        <LineChart
            className="mt-5 h-72"
            data={chartData}
            index="date"
            categories={['value']}
            showLegend={false}
            yAxisWidth={60}
            showXAxis
            showGridLines
            autoMinValue
            colors={['blue', 'indigo', 'rose-400', '#ffcc33']}
    />
      </div>
    </>
  );
}