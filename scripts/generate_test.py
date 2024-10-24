import random
import os

working_dir = os.path.dirname(os.path.dirname(os.path.realpath(__file__)))

SPECIAL_OP_PROB = 0.04

ops = ["B", "S"]
instruments = [
    "A1", "A2", "A3", "A4", "A5", "A6", "A7", "A8", "A9", "A10", "A11", "A12", "A13", "A14", "A15", "A16", "A17", "A18", "A19", "A20", "A21", "A22", "A23", "A24", "A25", "A26", "A27", "A28", "A29", "A30",
    "B1", "B2", "B3", "B4", "B5", "B6", "B7", "B8", "B9", "B10", "B11", "B12", "B13", "B14", "B15", "B16", "B17", "B18", "B19", "B20", "B21", "B22", "B23", "B24", "B25", "B26", "B27", "B28", "B29", "B30",
    "C1", "C2", "C3", "C4", "C5", "C6", "C7", "C8", "C9", "C10", "C11", "C12", "C13", "C14", "C15", "C16", "C17", "C18", "C19", "C20", "C21", "C22", "C23", "C24", "C25", "C26", "C27", "C28", "C29", "C30",
    "D1", "D2", "D3", "D4", "D5", "D6", "D7", "D8", "D9", "D10", "D11", "D12", "D13", "D14", "D15", "D16", "D17", "D18", "D19", "D20", "D21", "D22", "D23", "D24", "D25", "D26", "D27", "D28", "D29", "D30",
    "E1", "E2", "E3", "E4", "E5", "E6", "E7", "E8", "E9", "E10", "E11", "E12", "E13", "E14", "E15", "E16", "E17", "E18", "E19", "E20", "E21", "E22", "E23", "E24", "E25", "E26", "E27", "E28", "E29", "E30",
    "F1", "F2", "F3", "F4", "F5", "F6", "F7", "F8", "F9", "F10", "F11", "F12", "F13", "F14", "F15", "F16", "F17", "F18", "F19", "F20", "F21", "F22", "F23", "F24", "F25", "F26", "F27", "F28", "F29", "F30",
    "G1", "G2", "G3", "G4", "G5", "G6", "G7", "G8", "G9", "G10", "G11", "G12", "G13", "G14", "G15", "G16", "G17", "G18", "G19", "G20", "G21", "G22", "G23", "G24", "G25", "G26", "G27", "G28", "G29", "G30",
    "H1", "H2", "H3", "H4", "H5", "H6", "H7", "H8", "H9", "H10", "H11", "H12", "H13", "H14", "H15", "H16", "H17", "H18", "H19", "H20", "H21", "H22", "H23", "H24", "H25", "H26", "H27", "H28", "H29", "H30",
    "I1", "I2", "I3", "I4", "I5", "I6", "I7", "I8", "I9", "I10", "I11", "I12", "I13", "I14", "I15", "I16", "I17", "I18", "I19", "I20", "I21", "I22", "I23", "I24", "I25", "I26", "I27", "I28", "I29", "I30",
    "J1", "J2", "J3", "J4", "J5", "J6", "J7", "J8", "J9", "J10", "J11", "J12", "J13", "J14", "J15", "J16", "J17", "J18", "J19", "J20", "J21", "J22", "J23", "J24", "J25", "J26", "J27", "J28", "J29", "J30",
    "K1", "K2", "K3", "K4", "K5", "K6", "K7", "K8", "K9", "K10", "K11", "K12", "K13", "K14", "K15", "K16", "K17", "K18", "K19", "K20", "K21", "K22", "K23", "K24", "K25", "K26", "K27", "K28", "K29", "K30",
    "L1", "L2", "L3", "L4", "L5", "L6", "L7", "L8", "L9", "L10", "L11", "L12", "L13", "L14", "L15", "L16", "L17", "L18", "L19", "L20", "L21", "L22", "L23", "L24", "L25", "L26", "L27", "L28", "L29", "L30",
    "M1", "M2", "M3", "M4", "M5", "M6", "M7", "M8", "M9", "M10", "M11", "M12", "M13", "M14", "M15", "M16", "M17", "M18", "M19", "M20", "M21", "M22", "M23", "M24", "M25", "M26", "M27", "M28", "M29", "M30",
    "N1", "N2", "N3", "N4", "N5", "N6", "N7", "N8", "N9", "N10", "N11", "N12", "N13", "N14", "N15", "N16", "N17", "N18", "N19", "N20", "N21", "N22", "N23", "N24", "N25", "N26", "N27", "N28", "N29", "N30",
    "O1", "O2", "O3", "O4", "O5", "O6", "O7", "O8", "O9", "O10", "O11", "O12", "O13", "O14", "O15", "O16", "O17", "O18", "O19", "O20", "O21", "O22", "O23", "O24", "O25", "O26", "O27", "O28", "O29", "O30",
    "P1", "P2", "P3", "P4", "P5", "P6", "P7", "P8", "P9", "P10", "P11", "P12", "P13", "P14", "P15", "P16", "P17", "P18", "P19", "P20", "P21", "P22", "P23", "P24", "P25", "P26", "P27", "P28", "P29", "P30",
    "Q1", "Q2", "Q3", "Q4", "Q5", "Q6", "Q7", "Q8", "Q9", "Q10", "Q11", "Q12", "Q13", "Q14", "Q15", "Q16", "Q17", "Q18", "Q19", "Q20", "Q21", "Q22", "Q23", "Q24", "Q25", "Q26", "Q27", "Q28", "Q29", "Q30",
    "R1", "R2", "R3", "R4", "R5", "R6", "R7", "R8", "R9", "R10", "R11", "R12", "R13", "R14", "R15", "R16", "R17", "R18", "R19", "R20", "R21", "R22", "R23", "R24", "R25", "R26", "R27", "R28", "R29", "R30",
    "S1", "S2", "S3", "S4", "S5", "S6", "S7", "S8", "S9", "S10", "S11", "S12", "S13", "S14", "S15", "S16", "S17", "S18", "S19", "S20", "S21", "S22", "S23", "S24", "S25", "S26", "S27", "S28", "S29", "S30",
    "T1", "T2", "T3", "T4", "T5", "T6", "T7", "T8", "T9", "T10", "T11", "T12", "T13", "T14", "T15", "T16", "T17", "T18", "T19", "T20", "T21", "T22", "T23", "T24", "T25", "T26", "T27", "T28", "T29", "T30",
    "U1", "U2", "U3", "U4", "U5", "U6", "U7", "U8", "U9", "U10", "U11", "U12", "U13", "U14", "U15", "U16", "U17", "U18", "U19", "U20", "U21", "U22", "U23", "U24", "U25", "U26", "U27", "U28", "U29", "U30",
    "V1", "V2", "V3", "V4", "V5", "V6", "V7", "V8", "V9", "V10", "V11", "V12", "V13", "V14", "V15", "V16", "V17", "V18", "V19", "V20", "V21", "V22", "V23", "V24", "V25", "V26", "V27", "V28", "V29", "V30",
    "W1", "W2", "W3", "W4", "W5", "W6", "W7", "W8", "W9", "W10", "W11", "W12", "W13", "W14", "W15", "W16", "W17", "W18", "W19", "W20", "W21", "W22", "W23", "W24", "W25", "W26", "W27", "W28", "W29", "W30",
    "X1", "X2", "X3", "X4", "X5", "X6", "X7", "X8", "X9", "X10", "X11", "X12", "X13", "X14", "X15", "X16", "X17", "X18", "X19", "X20", "X21", "X22", "X23", "X24", "X25", "X26", "X27", "X28", "X29", "X30",
    "Y1", "Y2", "Y3", "Y4", "Y5", "Y6", "Y7", "Y8", "Y9", "Y10", "Y11", "Y12", "Y13", "Y14", "Y15", "Y16", "Y17", "Y18", "Y19", "Y20", "Y21", "Y22", "Y23", "Y24", "Y25", "Y26", "Y27", "Y28", "Y29", "Y30",
    "Z1", "Z2", "Z3", "Z4", "Z5", "Z6", "Z7", "Z8", "Z9", "Z10", "Z11", "Z12", "Z13", "Z14", "Z15", "Z16", "Z17", "Z18", "Z19", "Z20", "Z21", "Z22", "Z23", "Z24", "Z25", "Z26", "Z27", "Z28", "Z29", "Z30"
]

id = 0

def generate(
    num_clients: int, num_transactions: int, num_instruments: int,
    filename: str, has_cancel: bool = True, use_round_numbers: bool = False
):
    global id

    if num_instruments == -1:
        instrument_list = instruments
    else:
        instrument_list = instruments[:num_instruments]

    if num_clients < 1:
        print("Invalid number of clients")
        return 

    round_numbers = list(range(10, 101, 10))

    if num_clients == 1:
        print("Single client")
        with open(f"{working_dir}/scripts/{filename}", "w+") as f:
            f.write("1\n")
            f.write("o\n")

            for i in range(num_transactions):
                prob = random.random()
                if has_cancel and id > 0 and prob < 5 * SPECIAL_OP_PROB:
                    f.write(f"C {random.randint(0, id-1)}\n")
                else:
                    op = random.choice(ops)
                    instrument = random.choice(instrument_list)
                    count = random.choice(round_numbers) if use_round_numbers else random.randint(1, 100)
                    price = random.choice(round_numbers) if use_round_numbers else random.randint(1, 100)
                    f.write(f"{op} {id} {instrument} {count}\n")
                    id += 1
            f.write("x\n")
    
    else:

        if has_cancel:
            client_id_lookup = {} # map operation id to client id for cancel 

        print(f"{num_clients} clients")
        with open(filename, "w") as f:
            f.write(f"{num_clients}\n")
            f.write(f"0-{num_clients-1} o\n") # connect all threads to server

            for i in range(num_transactions):
                prob = random.random()

                num_sync = random.randint(2, num_clients)
                multi_clients = random.sample(range(num_clients), num_sync)
                multi_clients = ",".join([str(c) for c in multi_clients])

                single_client = random.randint(0, num_clients-1)

                if prob < SPECIAL_OP_PROB: # sleep
                    if random.random() < 0.5:
                        f.write(f"{multi_clients} s {random.randint(0, 1000)}\n") 
                    else:
                        f.write(f"{single_client} s {random.randint(0, 1000)}\n")
                elif SPECIAL_OP_PROB <= prob < 2 * SPECIAL_OP_PROB: # sync
                    f.write(f"{multi_clients} .\n")
                # elif id > 0 and 2 * SPECIAL_OP_PROB <= prob < 3 * SPECIAL_OP_PROB: # wait for active order to be matched
                #     if random.random() < 0.5:
                #         f.write(f"{multi_clients} w {random.randint(0, id-1)}\n")
                #     else:
                #         f.write(f"{single_client} w {random.randint(0, id-1)}\n")
                elif has_cancel and id > 0 and 3 * SPECIAL_OP_PROB <= prob < 4 * SPECIAL_OP_PROB: # cancel   
                    op_id = random.randint(0, id-1)
                    f.write(f"{client_id_lookup[op_id]} C {op_id}\n")
                else: # buy or sell
                    op = random.choice(ops)
                    instrument = random.choice(instrument_list)
                    count = random.choice(round_numbers) if use_round_numbers else random.randint(1, 100)
                    price = random.choice(round_numbers) if use_round_numbers else random.randint(1, 100)
                    f.write(f"{single_client} {op} {id} {instrument} {price} {count}\n")
                    if has_cancel:
                        client_id_lookup[id] = single_client
                    id += 1

            f.write(f"0-{num_clients-1} x\n") # disconnect all clients

        
if __name__=="__main__":
    num_clients = int(input("Number of clients: ") or 1)
    num_transactions = int(input("Number of test transactions: ") or 100)
    num_instruments = int(input(f"Number of instruments [1-{len(instruments)}]: ") or 1)
    filename = input("Name of the test file: ") or "test"
    filename += ".in"
    has_cancel = input("Include cancel operations? (Y/n) ") or "y"
    has_cancel = has_cancel.lower() == "y"
    use_round_numbers = input("Use round numbers for price and count? (Y/n) ") or "y"
    use_round_numbers = use_round_numbers.lower() == "y"
    generate(num_clients, num_transactions, num_instruments, filename, has_cancel, use_round_numbers)
