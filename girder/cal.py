G=86.24
P=98.01
GL=562.3
PL=639.4
GL0=569.2
PL0=646.8
AG=(GL+PL)/2
AP=(PL+PL0)/2

TG=100
TP=108.11
TGL=846.5
TPL=915.15
TGL0=840.0
TPL0=908.1
TAG=(TGL+TGL0)/2
TAP=(TPL+TPL0)/2

def M_printer(k: float):
    print("Cal for {}".format(k))
    print("GL  ", "{}*{}={}".format(k, GL, k*GL))
    print("PL  ", "{}*{}={}".format(k, PL, k*PL))
    print("GL0 ", "{}*{}={}".format(k, GL0, k*GL0))
    print("PL0 ", "{}*{}={}".format(k, PL0, k*PL0))
    print("AG  ", "{}*{}={}".format(k, AG, k*AG))
    print("AP  ", "{}*{}={}".format(k, AP, k*AP))
    print("")
    print("TGL ", "{}*{}={}".format(k, TGL, k*TGL))
    print("TPL ", "{}*{}={}".format(k, TPL, k*TPL))
    print("TGL0", "{}*{}={}".format(k, TGL0, k*TGL0))
    print("TPL0", "{}*{}={}".format(k, TPL0, k*TPL0))
    print("TAG ", "{}*{}={}".format(k, TAG, k*TAG))
    print("TAP ", "{}*{}={}".format(k, TAP, k*TAP))
    print("\n")

def V_cal(k: float):
    print("Cal for {}".format(k))
    print("{:2} {}*{}={}".format("G", k, G, k*G))
    print("{:2} {}*{}={}".format("P", k, P, k*P))
    print("")
    print("{:2} {}*{}={}".format("TG", k, TG, k*TG))
    print("{:2} {}*{}={}".format("TP", k, TP, k*TP))
    print("\n")

def cal():
    # M_printer(0.289)
    # M_printer(0.045)
    # M_printer(0.089)
    # M_printer(-0.133)
    # M_printer(0.2)
    # M_printer(0.229)
    # M_printer(0.125)
    # M_printer(-0.311)
    # M_printer(0.170)
    # M_printer(-0.089)

    # M_printer(0.059)
    # M_printer(0.096)
    # M_printer(0.03)
    V_cal(0.773)
    V_cal(-1.267)
    V_cal(1.0)
    V_cal(0.866)
    V_cal(-1.134)
    V_cal(-0.133)
    V_cal(1.0)
    V_cal(0.689)
    V_cal(-1.311)
    V_cal(1.222)
    V_cal(-0.089)
    V_cal(0.778)

cal()

