import type { Game } from "../types/app";

const CHARACTER_KEY_TO_ID: Record<string, string> = {
    group_voices: "00",
    kinzo: "01",
    krauss: "02",
    natsuhi: "03",
    jessica: "04",
    eva: "05",
    hideyoshi: "06",
    george: "07",
    rudolf: "08",
    kyrie: "09",
    battler: "10",
    ange: "11",
    rosa: "12",
    maria: "13",
    genji: "14",
    shannon: "15",
    kanon: "16",
    gohda: "17",
    kumasawa: "18",
    nanjo: "19",
    amakusa: "20",
    okonogi: "21",
    kasumi: "22",
    professor: "23",
    kawabata: "24",
    nanjo_son: "25",
    kumasawa_son: "26",
    beatrice: "27",
    bernkastel: "28",
    lambdadelta: "29",
    virgilia: "30",
    ronove: "31",
    gaap: "32",
    sakutarou: "33",
    eva_beatrice: "34",
    chiester_45: "35",
    chiester_410: "36",
    chiester_00: "37",
    lucifer: "38",
    leviathan: "39",
    satan: "40",
    belphegor: "41",
    mammon: "42",
    beelzebub: "43",
    asmodeus: "44",
    goat: "45",
    erika: "46",
    dlanor: "47",
    gertrude: "48",
    cornelia: "49",
    featherine: "50",
    zepar: "51",
    furfur: "52",
    lion: "53",
    will: "54",
    clair: "55",
    ikuko: "56",
    tohya: "57",
    kinzo_young: "58",
    bice: "59",
    beato_elder: "60",
    misc_voices: "99",
    narrator: "narrator",
};

const CHARACTER_ID_TO_KEY = Object.fromEntries(
    Object.entries(CHARACTER_KEY_TO_ID).map(([key, id]) => [id, key]),
) as Record<string, string>;

export function normalizeCharacterKey(value: string, game?: Game): string {
    if (!value || game === "higurashi") {
        return value;
    }
    return CHARACTER_ID_TO_KEY[value] || value;
}

export function toCharacterId(value: string, game?: Game): string {
    if (!value || game === "higurashi") {
        return value;
    }
    return CHARACTER_KEY_TO_ID[value] || value;
}
